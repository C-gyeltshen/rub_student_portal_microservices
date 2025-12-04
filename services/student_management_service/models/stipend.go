package models

import (
	"time"

	"gorm.io/gorm"
)

// StipendAllocation records stipend allocations for students
type StipendAllocation struct {
	gorm.Model
	AllocationID   string  `json:"allocation_id" gorm:"uniqueIndex;not null"`
	StudentID      uint    `json:"student_id" gorm:"not null"`
	Amount         float64 `json:"amount" gorm:"not null"`
	AllocationDate string  `json:"allocation_date"`
	Status         string  `json:"status" gorm:"default:'pending'"` // pending, approved, rejected, disbursed
	ApprovedBy     uint    `json:"approved_by"`
	ApprovalDate   string  `json:"approval_date"`
	Semester       int     `json:"semester"`
	AcademicYear   string  `json:"academic_year"`
	Remarks        string  `json:"remarks"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// StipendHistory tracks stipend payment history
type StipendHistory struct {
	gorm.Model
	TransactionID  string  `json:"transaction_id" gorm:"uniqueIndex;not null"`
	StudentID      uint    `json:"student_id" gorm:"not null"`
	AllocationID   uint    `json:"allocation_id"`
	Amount         float64 `json:"amount" gorm:"not null"`
	PaymentDate    string  `json:"payment_date"`
	TransactionStatus string `json:"transaction_status"` // success, failed, pending
	PaymentMethod  string  `json:"payment_method"` // bank_transfer, cash, etc.
	BankReference  string  `json:"bank_reference"`
	Remarks        string  `json:"remarks"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// StipendEligibility represents eligibility criteria and status
type StipendEligibility struct {
	StudentID        uint    `json:"student_id"`
	IsEligible       bool    `json:"is_eligible"`
	Reasons          []string `json:"reasons"`
	ExpectedAmount   float64 `json:"expected_amount"`
	AcademicStanding string  `json:"academic_standing"`
	AttendanceRate   float64 `json:"attendance_rate"`
	HasPendingIssues bool    `json:"has_pending_issues"`
}
