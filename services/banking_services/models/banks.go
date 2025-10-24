package models

import (
	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	Name string `json:"name"`
	
	// Optional: one-to-many relationship (a bank can have many student accounts)
	StudentBankDetails []StudentBankDetails `gorm:"foreignKey:BankID" json:"student_bank_details,omitempty"`
}
