package utils

import (
	"encoding/json"
	"net/http"
	"student_management_service/database"
	"student_management_service/middleware"
	"student_management_service/models"
	"time"
)

// LogAudit creates an audit log entry
func LogAudit(r *http.Request, entityType string, entityID uint, action string, changes interface{}) error {
	userID := middleware.GetUserIDFromContext(r.Context())
	
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return err
	}

	log := models.AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		UserID:     parseUserID(userID),
		Changes:    string(changesJSON),
		IPAddress:  getIPAddress(r),
		UserAgent:  r.UserAgent(),
		Timestamp:  time.Now(),
	}

	return database.DB.Create(&log).Error
}

// parseUserID converts string user ID to uint
func parseUserID(userID string) uint {
	// In production, you'd parse this properly
	// For now, returning 0 if invalid
	var id uint
	// Add proper parsing logic here
	return id
}

// getIPAddress extracts IP address from request
func getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	
	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	
	// Fall back to RemoteAddr
	return r.RemoteAddr
}
