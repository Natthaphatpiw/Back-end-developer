package middleware

import (
	"time"

	"github.com/Natthaphatpiw/Backend-with-GO-GIN/config"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/models"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		var token models.Token
		if err := config.DB.Where("token = ?", tokenString).First(&token).Error; err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.ExpiresAt.Before(time.Now()) {
			c.JSON(401, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		if !token.IsValid() {
			config.DB.Delete(&token)
			c.JSON(401, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Set("staff_id", token.StaffID)
		c.Set("hospital_id", token.HospitalID)

		c.Next()
	}
}
