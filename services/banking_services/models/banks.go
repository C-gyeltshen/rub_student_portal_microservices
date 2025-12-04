package models

import (
	"time"

	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type Bank struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Optional: one-to-many relationship (a bank can have many student accounts)
	StudentBankDetails []StudentBankDetails `gorm:"foreignKey:BankID" json:"student_bank_details,omitempty"`
}
