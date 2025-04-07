package models

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	FirstNameTh  string    `json:"first_name_th"`
	MiddleNameTh string    `json:"middle_name_th"`
	LastNameTh   string    `json:"last_name_th"`
	FirstNameEn  string    `json:"first_name_en"`
	MiddleNameEn string    `json:"middle_name_en"`
	LastNameEn   string    `json:"last_name_en"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	PatientHN    string    `json:"patient_hn"`
	NationalID   string    `json:"national_id"`
	PassportID   string    `json:"passport_id"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"` // M หรือ F
	HospitalID   uint      `json:"hospital_id"`
	Hospital     Hospital  `json:"hospital"`
}

type PatientSearchRequest struct {
	NationalID  string     `json:"national_id" form:"national_id" validate:"omitempty,len=13"`
	PassportID  string     `json:"passport_id" form:"passport_id" validate:"omitempty,min=5,max=20"`
	FirstName   string     `json:"first_name" form:"first_name" validate:"omitempty,min=2,max=50"`
	MiddleName  string     `json:"middle_name" form:"middle_name"`
	LastName    string     `json:"last_name" form:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth" form:"date_of_birth"`
	PhoneNumber string     `json:"phone_number" form:"phone_number"`
	Email       string     `json:"email" form:"email" validate:"omitempty,email"`
}

type PatientResponse struct {
	FirstNameTh  string    `json:"first_name_th"`
	MiddleNameTh string    `json:"middle_name_th"`
	LastNameTh   string    `json:"last_name_th"`
	FirstNameEn  string    `json:"first_name_en"`
	MiddleNameEn string    `json:"middle_name_en"`
	LastNameEn   string    `json:"last_name_en"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	PatientHN    string    `json:"patient_hn"`
	NationalID   string    `json:"national_id"`
	PassportID   string    `json:"passport_id"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"` // M หรือ F
}
