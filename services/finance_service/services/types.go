package services

import (
	"time"

	"github.com/google/uuid"
	"finance_service/models"
)

// Stipend represents a stipend (matches models.Stipend for gRPC compatibility)
type Stipend struct {
	ID            uuid.UUID
	StudentID     uuid.UUID
	Amount        float64
	StipendType   string
	PaymentDate   *time.Time
	PaymentStatus string
	PaymentMethod string
	JournalNumber string
	Notes         string
	CreatedAt     time.Time
	ModifiedAt    time.Time
}

// DeductionRule represents a deduction rule (matches models.DeductionRule for gRPC compatibility)
type DeductionRule struct {
	ID                        uuid.UUID
	RuleName                  string
	DeductionType             string
	DefaultAmount             float64
	MinAmount                 float64
	MaxAmount                 float64
	IsApplicableToFullScholar bool
	IsApplicableToSelfFunded  bool
	IsActive                  bool
	AppliesMonthly            bool
	AppliesAnnually           bool
	IsOptional                bool
	Priority                  int
	Description               string
	CreatedAt                 time.Time
	ModifiedAt                time.Time
}

// Deduction represents a deduction (matches models.Deduction for gRPC compatibility)
type Deduction struct {
	ID              uuid.UUID
	StudentID       uuid.UUID
	DeductionRuleID uuid.UUID
	StipendID       uuid.UUID
	Amount          float64
	DeductionType   string
	Description     string
	DeductionDate   time.Time
	ProcessingStatus string
	ApprovedBy      *uuid.UUID
	ApprovalDate    *time.Time
	RejectionReason string
	TransactionID   *uuid.UUID
	CreatedAt       time.Time
	ModifiedAt      time.Time
}

// convertModelStipendToService converts a models.Stipend to a services.Stipend
func convertModelStipendToService(ms *models.Stipend) *Stipend {
	return &Stipend{
		ID:            ms.ID,
		StudentID:     ms.StudentID,
		Amount:        ms.Amount,
		StipendType:   ms.StipendType,
		PaymentDate:   ms.PaymentDate,
		PaymentStatus: ms.PaymentStatus,
		PaymentMethod: ms.PaymentMethod,
		JournalNumber: ms.JournalNumber,
		Notes:         ms.Notes,
		CreatedAt:     ms.CreatedAt,
		ModifiedAt:    ms.ModifiedAt,
	}
}

// convertModelDeductionToService converts a models.Deduction to a services.Deduction
func convertModelDeductionToService(md *models.Deduction) *Deduction {
	return &Deduction{
		ID:              md.ID,
		StudentID:       md.StudentID,
		DeductionRuleID: md.DeductionRuleID,
		StipendID:       md.StipendID,
		Amount:          md.Amount,
		DeductionType:   md.DeductionType,
		Description:     md.Description,
		DeductionDate:   md.DeductionDate,
		ProcessingStatus: md.ProcessingStatus,
		ApprovedBy:      md.ApprovedBy,
		ApprovalDate:    md.ApprovalDate,
		RejectionReason: md.RejectionReason,
		TransactionID:   md.TransactionID,
		CreatedAt:       md.CreatedAt,
		ModifiedAt:      md.ModifiedAt,
	}
}

// convertModelDeductionRuleToService converts a models.DeductionRule to a services.DeductionRule
func convertModelDeductionRuleToService(mdr *models.DeductionRule) *DeductionRule {
	return &DeductionRule{
		ID:                        mdr.ID,
		RuleName:                  mdr.RuleName,
		DeductionType:             mdr.DeductionType,
		DefaultAmount:             mdr.BaseAmount,
		MinAmount:                 mdr.MinDeductionAmount,
		MaxAmount:                 mdr.MaxDeductionAmount,
		IsApplicableToFullScholar: mdr.IsApplicableToFullScholar,
		IsApplicableToSelfFunded:  mdr.IsApplicableToSelfFunded,
		IsActive:                  mdr.IsActive,
		AppliesMonthly:            mdr.AppliesMonthly,
		AppliesAnnually:           mdr.AppliesAnnually,
		IsOptional:                mdr.IsOptional,
		Priority:                  mdr.Priority,
		Description:               mdr.Description,
		CreatedAt:                 mdr.CreatedAt,
		ModifiedAt:                mdr.ModifiedAt,
	}
}
