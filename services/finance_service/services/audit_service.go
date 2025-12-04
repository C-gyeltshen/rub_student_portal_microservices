package services

import (
	"finance_service/models"
	"fmt"

	"gorm.io/gorm"
)

// AuditService handles all audit logging operations
type AuditService struct {
	db *gorm.DB
}

// NewAuditService creates a new audit service
func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

// LogAction logs a financial action to the audit trail
func (as *AuditService) LogAction(action, entityType, entityID, financeOfficer, description string) error {
	auditLog := models.NewAuditLog(action, entityType, entityID, financeOfficer, description)
	if err := as.db.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to log audit action: %w", err)
	}
	return nil
}

// LogActionWithValues logs an action with old and new values
func (as *AuditService) LogActionWithValues(action, entityType, entityID, financeOfficer, description string, oldValues, newValues interface{}) error {
	auditLog := models.NewAuditLog(action, entityType, entityID, financeOfficer, description)
	
	if err := auditLog.SetOldValues(oldValues); err != nil {
		return fmt.Errorf("failed to set old values: %w", err)
	}
	
	if err := auditLog.SetNewValues(newValues); err != nil {
		return fmt.Errorf("failed to set new values: %w", err)
	}
	
	if err := as.db.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to log audit action with values: %w", err)
	}
	return nil
}

// LogError logs a failed action
func (as *AuditService) LogError(action, entityType, entityID, financeOfficer, description, errorMsg string) error {
	auditLog := models.NewAuditLog(action, entityType, entityID, financeOfficer, description)
	auditLog.SetError(errorMsg)
	
	if err := as.db.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to log audit error: %w", err)
	}
	return nil
}

// GetAuditLogs retrieves audit logs with optional filtering
func (as *AuditService) GetAuditLogs(limit, offset int, filters map[string]interface{}) ([]models.AuditLog, int64, error) {
	var auditLogs []models.AuditLog
	var total int64
	
	query := as.db
	
	// Apply filters
	if action, ok := filters["action"]; ok {
		query = query.Where("action = ?", action)
	}
	if entityType, ok := filters["entity_type"]; ok {
		query = query.Where("entity_type = ?", entityType)
	}
	if financeOfficer, ok := filters["finance_officer"]; ok {
		query = query.Where("finance_officer = ?", financeOfficer)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if startDate, ok := filters["start_date"]; ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"]; ok {
		query = query.Where("created_at <= ?", endDate)
	}
	
	// Get total count
	if err := query.Model(&models.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}
	
	// Get paginated results
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&auditLogs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch audit logs: %w", err)
	}
	
	return auditLogs, total, nil
}

// GetAuditLogsByEntity retrieves all audit logs for a specific entity
func (as *AuditService) GetAuditLogsByEntity(entityType, entityID string) ([]models.AuditLog, error) {
	var auditLogs []models.AuditLog
	if err := as.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&auditLogs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch audit logs for entity: %w", err)
	}
	return auditLogs, nil
}

// GetAuditLogsByOfficer retrieves all audit logs for a specific finance officer
func (as *AuditService) GetAuditLogsByOfficer(financeOfficer string, limit, offset int) ([]models.AuditLog, int64, error) {
	var auditLogs []models.AuditLog
	var total int64
	
	if err := as.db.Model(&models.AuditLog{}).
		Where("finance_officer = ?", financeOfficer).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count officer logs: %w", err)
	}
	
	if err := as.db.Where("finance_officer = ?", financeOfficer).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch officer audit logs: %w", err)
	}
	
	return auditLogs, total, nil
}
