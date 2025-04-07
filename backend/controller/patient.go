package controller

import (
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/config"
	"github.com/Natthaphatpiw/Backend-with-GO-GIN/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPatient(c *gin.Context) {
	id := c.Param("id")

	var patient models.Patient
	if err := config.DB.
		Where("national_id = ? OR passport_id = ?", id, id).
		First(&patient).Error; err != nil {
		c.JSON(404, gin.H{"error": "Patient not found"})
		return
	}

	response := models.PatientResponse{
		FirstNameTh:  patient.FirstNameTh,
		MiddleNameTh: patient.MiddleNameTh,
		LastNameTh:   patient.LastNameTh,
		FirstNameEn:  patient.FirstNameEn,
		MiddleNameEn: patient.MiddleNameEn,
		LastNameEn:   patient.LastNameEn,
		DateOfBirth:  patient.DateOfBirth,
		PatientHN:    patient.PatientHN,
		NationalID:   patient.NationalID,
		PassportID:   patient.PassportID,
		PhoneNumber:  patient.PhoneNumber,
		Email:        patient.Email,
		Gender:       patient.Gender,
	}

	c.JSON(200, response)
}

func SearchPatients(c *gin.Context) {
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		c.JSON(500, gin.H{"error": "Hospital ID not found in context"})
		return
	}

	var searchRequest models.PatientSearchRequest
	if err := c.ShouldBindQuery(&searchRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := config.DB.Model(&models.Patient{}).Where("hospital_id = ?", hospitalID)

	if searchRequest.NationalID != "" {
		query = query.Where("national_id LIKE ?", "%"+searchRequest.NationalID+"%")
	}
	if searchRequest.PassportID != "" {
		query = query.Where("passport_id LIKE ?", "%"+searchRequest.PassportID+"%")
	}
	if searchRequest.FirstName != "" {
		query = query.Where("first_name_th LIKE ? OR first_name_en LIKE ?",
			"%"+searchRequest.FirstName+"%", "%"+searchRequest.FirstName+"%")
	}
	if searchRequest.MiddleName != "" {
		query = query.Where("middle_name_th LIKE ? OR middle_name_en LIKE ?",
			"%"+searchRequest.MiddleName+"%", "%"+searchRequest.MiddleName+"%")
	}
	if searchRequest.LastName != "" {
		query = query.Where("last_name_th LIKE ? OR last_name_en LIKE ?",
			"%"+searchRequest.LastName+"%", "%"+searchRequest.LastName+"%")
	}
	if searchRequest.DateOfBirth != nil {
		query = query.Where("date_of_birth = ?", searchRequest.DateOfBirth)
	}
	if searchRequest.PhoneNumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+searchRequest.PhoneNumber+"%")
	}
	if searchRequest.Email != "" {
		query = query.Where("email LIKE ?", "%"+searchRequest.Email+"%")
	}

	// Execute query
	var patients []models.Patient
	if err := query.Find(&patients).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(200, gin.H{"data": []models.PatientResponse{}})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to search patients"})
		return
	}

	var responses []models.PatientResponse
	for _, patient := range patients {
		responses = append(responses, models.PatientResponse{
			FirstNameTh:  patient.FirstNameTh,
			MiddleNameTh: patient.MiddleNameTh,
			LastNameTh:   patient.LastNameTh,
			FirstNameEn:  patient.FirstNameEn,
			MiddleNameEn: patient.MiddleNameEn,
			LastNameEn:   patient.LastNameEn,
			DateOfBirth:  patient.DateOfBirth,
			PatientHN:    patient.PatientHN,
			NationalID:   patient.NationalID,
			PassportID:   patient.PassportID,
			PhoneNumber:  patient.PhoneNumber,
			Email:        patient.Email,
			Gender:       patient.Gender,
		})
	}

	c.JSON(200, gin.H{"data": responses})
}
