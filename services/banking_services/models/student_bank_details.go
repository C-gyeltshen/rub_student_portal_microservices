package models

import (
	"time"

	"gorm.io/gorm"
)

type StudentBankDetails struct {
	gorm.Model
	StudentID         int            `json:"student_id"`
	BankID            uint           `json:"bank_id"` // Must be uint to match Bank.ID type
	AccountNumber     string         `json:"account_number"`
	AccountHolderName string         `json:"account_holder_name"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	// Foreign key relation
	Bank Bank `gorm:"foreignKey:BankID;references:ID" json:"bank,omitempty"`
}