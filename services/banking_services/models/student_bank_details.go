package models

import (
	"time"

	"gorm.io/gorm"
)

type StudentBankDetails struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StudentID         string    `gorm:"type:uuid" json:"student_id"`
	BankID            string    `gorm:"type:uuid" json:"bank_id"`
	AccountNumber     string    `json:"account_number"`
	AccountHolderName string    `json:"account_holder_name"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	// Foreign key relation
	Bank Bank `gorm:"foreignKey:BankID;references:ID" json:"bank,omitempty"`
}
