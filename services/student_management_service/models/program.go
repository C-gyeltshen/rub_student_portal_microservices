package models

import (
	"time"

	"gorm.io/gorm"
)

// College represents an academic college/institution
type College struct {
	gorm.Model
	Code        string `json:"code" gorm:"uniqueIndex;not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Location    string `json:"location"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	AllowSelfFinancedStipend bool `json:"allow_self_financed_stipend" gorm:"default:false"` // Whether college allows self-financed students to get stipend
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Program represents an academic program (e.g., BSc IT, BEd)
type Program struct {
	gorm.Model
	Code             string  `json:"code" gorm:"uniqueIndex;not null"`
	Name             string  `json:"name" gorm:"not null"`
	Description      string  `json:"description"`
	Level            string  `json:"level"` // undergraduate, postgraduate
	DurationYears    int     `json:"duration_years"`
	DurationSemesters int    `json:"duration_semesters"`
	CollegeID        uint    `json:"college_id"`
	College          College `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
	
	// Stipend Information
	HasStipend       bool    `json:"has_stipend" gorm:"default:false"`
	StipendAmount    float64 `json:"stipend_amount"`
	StipendType      string  `json:"stipend_type"` // monthly, semester, annual
	
	IsActive         bool    `json:"is_active" gorm:"default:true"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
