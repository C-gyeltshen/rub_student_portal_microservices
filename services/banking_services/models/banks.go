package models

import (
	"time"

	"gorm.io/gorm"
)

type Bank struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name               string    `gorm:"index:,type:btree" json:"name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Optional: one-to-many relationship (a bank can have many student accounts)
	StudentBankDetails []StudentBankDetails `gorm:"foreignKey:BankID" json:"student_bank_details,omitempty"`
}
