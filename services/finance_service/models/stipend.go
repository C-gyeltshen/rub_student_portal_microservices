package models

import (
	"github.com/google/uuid"
	"time"
)

// Stipend represents a stipend payment record for a student
type Stipend struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	StudentID     uuid.UUID  `gorm:"type:uuid;not null" json:"student_id"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	StipendType   string     `gorm:"type:varchar(50);not null" json:"stipend_type"` // full-scholarship, self-funded
	PaymentDate   *time.Time `gorm:"type:timestamptz" json:"payment_date"`
	PaymentStatus string     `gorm:"type:varchar(50);default:'Pending'" json:"payment_status"` // Pending, Processed, Failed
	PaymentMethod string     `gorm:"type:varchar(50)" json:"payment_method"`                   // Bank_transfer, E-payment
	JournalNumber string     `gorm:"type:text;unique;not null" json:"journal_number"`
	TransactionID *uuid.UUID `gorm:"type:uuid" json:"transaction_id"`
	Notes         string     `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	ModifiedAt    time.Time  `gorm:"type:timestamptz;autoUpdateTime" json:"modified_at"`
}

// TableName specifies the table name for Stipend model
func (Stipend) TableName() string {
	return "stipends"
}
