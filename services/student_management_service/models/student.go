package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	
	// Link to User Service
	UserID      uint   `json:"user_id" gorm:"uniqueIndex;not null"` // Links to User_Service
	
	// Student Information
	StudentID   string `json:"student_id" gorm:"uniqueIndex;not null"`
	FirstName   string `json:"first_name" gorm:"not null"`
	LastName    string `json:"last_name" gorm:"not null"`
	Email       string `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	CID         string `json:"cid" gorm:"uniqueIndex"` // Citizenship ID
	
	// Address Information
	PermanentAddress string `json:"permanent_address"`
	CurrentAddress   string `json:"current_address"`
	
	// Academic Information
	ProgramID        uint    `json:"program_id"`
	Program          Program `gorm:"foreignKey:ProgramID" json:"program,omitempty"`
	CollegeID        uint    `json:"college_id"`
	College          College `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
	YearOfStudy      int     `json:"year_of_study"`
	Semester         int     `json:"semester"`
	EnrollmentDate   string  `json:"enrollment_date"`
	GraduationDate   string  `json:"graduation_date"`
	Status           string `json:"status" gorm:"default:'active'"` // active, inactive, graduated, suspended
	AcademicStanding string `json:"academic_standing"` // good, probation, etc.
	FinancingType    string `json:"financing_type"` // scholarship, self-financed
	
	// Guardian Information
	GuardianName        string `json:"guardian_name"`
	GuardianPhoneNumber string `json:"guardian_phone_number"`
	GuardianRelation    string `json:"guardian_relation"`
	
	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
