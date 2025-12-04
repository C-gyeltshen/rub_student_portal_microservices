package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"finance_service/database"
	"finance_service/services"

	"github.com/go-chi/chi/v5"
)

// AuditHandler handles audit log requests
type AuditHandler struct {
	auditService *services.AuditService
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler() *AuditHandler {
	return &AuditHandler{
		auditService: services.NewAuditService(database.DB),
	}
}

// GetAuditLogs returns all audit logs with optional filters and pagination
func (ah *AuditHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	entityType := r.URL.Query().Get("entity_type")
	financeOfficer := r.URL.Query().Get("finance_officer")
	status := r.URL.Query().Get("status")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Parse pagination
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	// Parse dates
	filters := make(map[string]interface{})
	if action != "" {
		filters["action"] = action
	}
	if entityType != "" {
		filters["entity_type"] = entityType
	}
	if financeOfficer != "" {
		filters["finance_officer"] = financeOfficer
	}
	if status != "" {
		filters["status"] = status
	}
	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filters["start_date"] = t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filters["end_date"] = t
		}
	}

	logs, total, err := ah.auditService.GetAuditLogs(limit, offset, filters)
	if err != nil {
		log.Printf("Error fetching audit logs: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch audit logs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       logs,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"count":      len(logs),
	})
}

// GetAuditLogsByEntity returns audit logs for a specific entity
func (ah *AuditHandler) GetAuditLogsByEntity(w http.ResponseWriter, r *http.Request) {
	entityType := chi.URLParam(r, "entity_type")
	entityID := chi.URLParam(r, "entity_id")

	if entityType == "" || entityID == "" {
		http.Error(w, "Missing entity_type or entity_id parameter", http.StatusBadRequest)
		return
	}

	logs, err := ah.auditService.GetAuditLogsByEntity(entityType, entityID)
	if err != nil {
		log.Printf("Error fetching audit logs for entity: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch audit logs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
		"logs":        logs,
		"count":       len(logs),
	})
}

// GetAuditLogsByOfficer returns audit logs for a specific finance officer
func (ah *AuditHandler) GetAuditLogsByOfficer(w http.ResponseWriter, r *http.Request) {
	financeOfficer := chi.URLParam(r, "officer")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if financeOfficer == "" {
		http.Error(w, "Missing finance officer parameter", http.StatusBadRequest)
		return
	}

	// Parse pagination
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	logs, total, err := ah.auditService.GetAuditLogsByOfficer(financeOfficer, limit, offset)
	if err != nil {
		log.Printf("Error fetching audit logs for officer: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch audit logs: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"finance_officer": financeOfficer,
		"data":            logs,
		"total":           total,
		"limit":           limit,
		"offset":          offset,
		"count":           len(logs),
	})
}
