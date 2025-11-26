package models

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog records all changes made to student data for auditing
type AuditLog struct {
	gorm.Model
	EntityType   string `json:"entity_type"` // student, stipend, program, etc.
	EntityID     uint   `json:"entity_id"`
	Action       string `json:"action"` // create, update, delete
	UserID       uint   `json:"user_id"` // who made the change
	Changes      string `json:"changes" gorm:"type:jsonb"` // JSON of what changed
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	Timestamp    time.Time `json:"timestamp"`
	
	CreatedAt time.Time `json:"created_at"`
}
