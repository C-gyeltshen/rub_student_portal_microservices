package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentBankDetails struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	StudentID         uuid.UUID      `gorm:"type:uuid" json:"student_id"`
	BankID            uuid.UUID      `gorm:"type:uuid" json:"bank_id"`
	AccountNumber     string         `json:"account_number"`
	AccountHolderName string         `json:"account_holder_name"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign key relation
	Bank Bank `gorm:"foreignKey:BankID;references:ID" json:"bank,omitempty"`
}
