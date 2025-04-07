package controller

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/Natthaphatpiw/Backend-with-GO-GIN/config"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateStaff(c *gin.Context) {
	var request models.StaffCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	var hospital models.Hospital
	if err := tx.First(&hospital, request.HospitalID).Error; err != nil {
		tx.Rollback()
		c.JSON(404, gin.H{"error": "Hospital not found"})
		return
	}

	var existingStaff models.Staff
	if err := config.DB.Where("username = ?", request.Username).First(&existingStaff).Error; err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	staff := models.Staff{
		Username:   request.Username,
		Password:   string(hashedPassword),
		Name:       request.Name,
		Email:      request.Email,
		HospitalID: request.HospitalID,
	}

	if err := tx.Create(&staff).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to create staff"})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"data": staff.ToResponse()})
}

func LoginStaff(c *gin.Context) {
	var request models.StaffLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var staff models.Staff
	if err := config.DB.Where("username = ? AND hospital_id = ?", request.Username, request.HospitalID).First(&staff).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(request.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	tokenString := hex.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(24 * time.Hour)

	token := models.Token{
		Token:      tokenString,
		StaffID:    staff.ID,
		HospitalID: staff.HospitalID,
		ExpiresAt:  expiresAt,
	}
	if err := config.DB.Create(&token).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save token"})
		return
	}

	c.JSON(200, gin.H{
		"data": models.TokenResponse{
			Token:     tokenString,
			ExpiresAt: expiresAt,
			Staff:     staff.ToResponse(),
		},
	})
}
