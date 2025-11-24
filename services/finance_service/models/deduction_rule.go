package models

import (
	"time"
	"github.com/google/uuid"
)

// DeductionRule defines rules for deductions that can be applied to stipends
type DeductionRule struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	RuleName            string    `gorm:"type:varchar(100);not null;unique;index" json:"rule_name"`
	DeductionType       string    `gorm:"type:varchar(100);not null" json:"deduction_type"` // hostel, electricity, mess_fees, etc
	Description         string    `gorm:"type:text" json:"description"`
	BaseAmount          float64   `gorm:"type:decimal(10,2);not null;check:base_amount >= 0" json:"base_amount"`
	MaxDeductionAmount  float64   `gorm:"type:decimal(10,2);not null;check:max_deduction_amount >= 0" json:"max_deduction_amount"`
	MinDeductionAmount  float64   `gorm:"type:decimal(10,2);default:0;check:min_deduction_amount >= 0" json:"min_deduction_amount"`
	IsApplicableToFullScholar  bool   `gorm:"type:boolean;default:false" json:"is_applicable_to_full_scholar"`
	IsApplicableToSelfFunded   bool   `gorm:"type:boolean;default:true" json:"is_applicable_to_self_funded"`
	IsActive            bool      `gorm:"type:boolean;default:true;index" json:"is_active"`
	AppliesMonthly      bool      `gorm:"type:boolean;default:false" json:"applies_monthly"`
	AppliesAnnually     bool      `gorm:"type:boolean;default:false" json:"applies_annually"`
	IsOptional          bool      `gorm:"type:boolean;default:false" json:"is_optional"` // true if deduction is optional, false if mandatory
	Priority            int       `gorm:"type:integer;default:0" json:"priority"` // Higher priority deductions applied first
	CreatedBy           *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt           time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	ModifiedBy          *uuid.UUID `gorm:"type:uuid" json:"modified_by"`
	ModifiedAt          time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"modified_at"`
}

// TableName specifies the table name for DeductionRule model
func (DeductionRule) TableName() string {
	return "deduction_rules"
}
