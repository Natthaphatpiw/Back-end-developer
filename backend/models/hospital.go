package models

import (
	"gorm.io/gorm"
)

type Hospital struct {
	gorm.Model
	Name     string
	Location string

	Staffs   []Staff
	Patients []Patient
}
