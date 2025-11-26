package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Transaction represents a money transfer transaction for stipend distribution
type Transaction struct {
	ID                uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	StipendID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"stipend_id"`
	StudentID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"student_id"`
	Amount            float64       `gorm:"type:decimal(10,2);not null;check:amount > 0" json:"amount"`
	SourceAccount     string        `gorm:"type:varchar(255)" json:"source_account"`         // College/Institution account
	DestinationAccount string       `gorm:"type:varchar(255);not null" json:"destination_account"` // Student's account
	DestinationBank   string        `gorm:"type:varchar(255)" json:"destination_bank"`
	TransactionType   string        `gorm:"type:varchar(50);default:'STIPEND'" json:"transaction_type"` // STIPEND, REFUND, etc
	Status            string        `gorm:"type:varchar(50);default:'PENDING'" json:"status"` // PENDING, PROCESSING, SUCCESS, FAILED, CANCELLED
	PaymentMethod     string        `gorm:"type:varchar(50)" json:"payment_method"` // BANK_TRANSFER, E_PAYMENT, etc
	ReferenceNumber   sql.NullString `gorm:"type:varchar(255);uniqueIndex:,where:reference_number IS NOT NULL" json:"reference_number"` // Unique reference from payment gateway (NULL values allowed)
	ErrorMessage      string        `gorm:"type:text" json:"error_message"` // Error details if transaction failed
	Remarks           string        `gorm:"type:text" json:"remarks"`
	InitiatedAt       time.Time     `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP" json:"initiated_at"`
	ProcessedAt       *time.Time    `gorm:"type:timestamptz" json:"processed_at"`
	CompletedAt       *time.Time    `gorm:"type:timestamptz" json:"completed_at"`
	CreatedAt         time.Time     `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	ModifiedAt        time.Time     `gorm:"type:timestamptz;autoUpdateTime" json:"modified_at"`

	// Relationships
	Stipend Stipend `gorm:"foreignKey:StipendID;references:ID" json:"stipend,omitempty"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// TransactionStatus constants
const (
	TransactionStatusPending    = "PENDING"
	TransactionStatusProcessing = "PROCESSING"
	TransactionStatusSuccess    = "SUCCESS"
	TransactionStatusFailed     = "FAILED"
	TransactionStatusCancelled  = "CANCELLED"
)

// TransactionType constants
const (
	TransactionTypeStipend = "STIPEND"
	TransactionTypeRefund  = "REFUND"
)

// PaymentMethods constants
const (
	PaymentMethodBankTransfer = "BANK_TRANSFER"
	PaymentMethodEPayment     = "E_PAYMENT"
)
