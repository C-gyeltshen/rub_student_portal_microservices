package models

import (
	"github.com/google/uuid"
	"time"
)

// Deduction represents a deduction applied to a student's stipend
type Deduction struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	StudentID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	DeductionRuleID  uuid.UUID  `gorm:"type:uuid;not null" json:"deduction_rule_id"`
	StipendID        uuid.UUID  `gorm:"type:uuid;not null" json:"stipend_id"`
	Amount           float64    `gorm:"type:decimal(10,2);not null;check:amount >= 0" json:"amount"`
	DeductionType    string     `gorm:"type:varchar(100);not null" json:"deduction_type"` // hostel, electricity, mess_fees, etc
	Description      string     `gorm:"type:text" json:"description"`
	DeductionDate    time.Time  `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"deduction_date"`
	ProcessingStatus string     `gorm:"type:varchar(50);default:'Pending'" json:"processing_status"` // Pending, Approved, Processed, Rejected
	ApprovedBy       *uuid.UUID `gorm:"type:uuid" json:"approved_by"`
	ApprovalDate     *time.Time `gorm:"type:timestamptz" json:"approval_date"`
	RejectionReason  string     `gorm:"type:text" json:"rejection_reason"`
	TransactionID    *uuid.UUID `gorm:"type:uuid" json:"transaction_id"`
	CreatedAt        time.Time  `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	ModifiedAt       time.Time  `gorm:"type:timestamptz;autoUpdateTime" json:"modified_at"`
}

// TableName specifies the table name for Deduction model
func (Deduction) TableName() string {
	return "deductions"
}
