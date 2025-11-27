package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name            string `json:"last_name" gorm:"not null"`
	RubIDCardNumber string `json:"rub_id_card_number" gorm:"uniqueIndex;not null"`
	Email               string `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber         string `json:"phone_number"`
	DateOfBirth         string `json:"date_of_birth"`

	ProgramID           uint   `json:"program_id"`
	Program             Program `gorm:"foreignKey:ProgramID" json:"program,omitempty"`
	CollegeID           uint    `json:"college_id"`
	College             College `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
	UserID              uint    `json:"user_id"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
