package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Token      string    `json:"token" gorm:"uniqueIndex"`
	StaffID    uint      `json:"staff_id"`
	Staff      Staff     `json:"-"`
	HospitalID uint      `json:"hospital_id"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type TokenResponse struct {
	Token     string        `json:"token"`
	ExpiresAt time.Time     `json:"expires_at"`
	Staff     StaffResponse `json:"staff"`
}

func (t *Token) IsValid() bool {
	return time.Now().Before(t.ExpiresAt)
}
