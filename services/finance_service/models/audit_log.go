package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditLog tracks all financial operations for compliance and transparency
type AuditLog struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Action        string    `gorm:"index;type:varchar(50)" json:"action"` // CREATE, UPDATE, DELETE, VIEW
	EntityType    string    `gorm:"index;type:varchar(50)" json:"entity_type"` // STIPEND, DEDUCTION_RULE, TRANSACTION
	EntityID      string    `gorm:"index;type:uuid" json:"entity_id"`     // ID of the entity being acted upon
	FinanceOfficer string   `gorm:"type:varchar(255)" json:"finance_officer"` // Email/ID of the officer
	Description   string    `gorm:"type:text" json:"description"`
	OldValues     string    `gorm:"type:jsonb;default:'{}'" json:"old_values,omitempty"`   // Previous data for updates
	NewValues     string    `gorm:"type:jsonb;default:'{}'" json:"new_values,omitempty"`   // New data
	Status        string    `gorm:"type:varchar(50);default:'SUCCESS'" json:"status"`      // SUCCESS, FAILED
	ErrorMessage  string    `gorm:"type:text" json:"error_message,omitempty"`
	IPAddress     string    `gorm:"type:varchar(50)" json:"ip_address,omitempty"`
	UserAgent     string    `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(action, entityType, entityID, financeOfficer, description string) *AuditLog {
	return &AuditLog{
		ID:             uuid.New().String(),
		Action:         action,
		EntityType:     entityType,
		EntityID:       entityID,
		FinanceOfficer: financeOfficer,
		Description:    description,
		Status:         "SUCCESS",
	}
}

// SetOldValues sets the old values for audit tracking
func (a *AuditLog) SetOldValues(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	a.OldValues = string(jsonData)
	return nil
}

// SetNewValues sets the new values for audit tracking
func (a *AuditLog) SetNewValues(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	a.NewValues = string(jsonData)
	return nil
}

// SetError marks the audit log as failed
func (a *AuditLog) SetError(errMsg string) {
	a.Status = "FAILED"
	a.ErrorMessage = errMsg
}

// SetIPAddress sets the IP address of the requester
func (a *AuditLog) SetIPAddress(ip string) {
	a.IPAddress = ip
}

// SetUserAgent sets the user agent
func (a *AuditLog) SetUserAgent(ua string) {
	a.UserAgent = ua
}
