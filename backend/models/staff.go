package models

import (
	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	Username   string `json:"username" gorm:"uniqueIndex"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	HospitalID uint   `json:"hospital_id"`
	Hospital   Hospital
}

type StaffLoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	HospitalID uint   `json:"hospital" binding:"required"`
}

type StaffCreateRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email"`
	HospitalID uint   `json:"hospital" binding:"required"`
}

type StaffResponse struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	HospitalID uint   `json:"hospital_id"`
}

func (s *Staff) ToResponse() StaffResponse {
	return StaffResponse{
		ID:         s.ID,
		Username:   s.Username,
		Name:       s.Name,
		Email:      s.Email,
		HospitalID: s.HospitalID,
	}
}
