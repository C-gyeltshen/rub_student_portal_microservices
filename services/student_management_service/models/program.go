package models

import (
	"time"

	"gorm.io/gorm"
)

// College represents an academic college/institution
type College struct {
	gorm.Model
	Name                     string         `json:"name" gorm:"not null"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Program represents an academic program (e.g., BSc IT, BEd)
type Program struct {
	gorm.Model
	Name          string  `json:"name" gorm:"not null"`
	CollegeID     uint    `json:"college_id"`
	College       College `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
