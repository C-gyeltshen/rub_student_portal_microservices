package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"finance_service/models"
	"finance_service/services"
)

// DeductionHandler handles HTTP requests for deduction operations
type DeductionHandler struct {
	stipendService *services.StipendService
}

// NewDeductionHandler creates a new deduction handler
func NewDeductionHandler() *DeductionHandler {
	return &DeductionHandler{
		stipendService: services.NewStipendService(),
	}
}

// CreateDeductionRuleRequest represents the request body for creating a deduction rule
type CreateDeductionRuleRequest struct {
	RuleName                  string  `json:"rule_name"`
	DeductionType             string  `json:"deduction_type"`
	Description               string  `json:"description"`
	BaseAmount                float64 `json:"base_amount"`
	MaxDeductionAmount        float64 `json:"max_deduction_amount"`
	MinDeductionAmount        float64 `json:"min_deduction_amount,omitempty"`
	IsApplicableToFullScholar bool    `json:"is_applicable_to_full_scholar"`
	IsApplicableToSelfFunded  bool    `json:"is_applicable_to_self_funded"`
	AppliesMonthly            bool    `json:"applies_monthly"`
	AppliesAnnually           bool    `json:"applies_annually"`
	IsOptional                bool    `json:"is_optional"`
	Priority                  int     `json:"priority"`
}

// DeductionRuleResponse represents the response for a deduction rule
type DeductionRuleResponse struct {
	ID                        string    `json:"id"`
	RuleName                  string    `json:"rule_name"`
	DeductionType             string    `json:"deduction_type"`
	Description               string    `json:"description"`
	BaseAmount                float64   `json:"base_amount"`
	MaxDeductionAmount        float64   `json:"max_deduction_amount"`
	MinDeductionAmount        float64   `json:"min_deduction_amount"`
	IsApplicableToFullScholar bool      `json:"is_applicable_to_full_scholar"`
	IsApplicableToSelfFunded  bool      `json:"is_applicable_to_self_funded"`
	IsActive                  bool      `json:"is_active"`
	AppliesMonthly            bool      `json:"applies_monthly"`
	AppliesAnnually           bool      `json:"applies_annually"`
	IsOptional                bool      `json:"is_optional"`
	Priority                  int       `json:"priority"`
	CreatedAt                 time.Time `json:"created_at"`
	ModifiedAt                time.Time `json:"modified_at"`
}

// DeductionResponse represents the response for a deduction record
type DeductionResponse struct {
	ID               string     `json:"id"`
	StudentID        string     `json:"student_id"`
	DeductionRuleID  string     `json:"deduction_rule_id"`
	StipendID        string     `json:"stipend_id"`
	Amount           float64    `json:"amount"`
	DeductionType    string     `json:"deduction_type"`
	Description      string     `json:"description"`
	DeductionDate    time.Time  `json:"deduction_date"`
	ProcessingStatus string     `json:"processing_status"`
	ApprovedBy       *string    `json:"approved_by"`
	ApprovalDate     *time.Time `json:"approval_date"`
	RejectionReason  string     `json:"rejection_reason"`
	CreatedAt        time.Time  `json:"created_at"`
	ModifiedAt       time.Time  `json:"modified_at"`
}

// CreateDeductionRule handles POST /api/deduction-rules - creates a new deduction rule
func (h *DeductionHandler) CreateDeductionRule(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDeductionRule handler called")

	var req CreateDeductionRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	rule := &models.DeductionRule{
		RuleName:                  req.RuleName,
		DeductionType:             req.DeductionType,
		Description:               req.Description,
		BaseAmount:                req.BaseAmount,
		MaxDeductionAmount:        req.MaxDeductionAmount,
		MinDeductionAmount:        req.MinDeductionAmount,
		IsApplicableToFullScholar: req.IsApplicableToFullScholar,
		IsApplicableToSelfFunded:  req.IsApplicableToSelfFunded,
		AppliesMonthly:            req.AppliesMonthly,
		AppliesAnnually:           req.AppliesAnnually,
		IsOptional:                req.IsOptional,
		Priority:                  req.Priority,
		IsActive:                  true,
	}

	if err := h.stipendService.CreateDeductionRule(rule); err != nil {
		log.Printf("Error creating deduction rule: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create deduction rule: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toDeductionRuleResponse(rule))
}

// GetDeductionRule handles GET /api/deduction-rules/{ruleID} - retrieves a deduction rule by ID
func (h *DeductionHandler) GetDeductionRule(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDeductionRule handler called")

	ruleIDStr := chi.URLParam(r, "ruleID")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		http.Error(w, "Invalid rule ID format", http.StatusBadRequest)
		return
	}

	rule, err := h.stipendService.GetDeductionRuleByID(ruleID)
	if err != nil {
		log.Printf("Error fetching deduction rule: %v", err)
		http.Error(w, "Deduction rule not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toDeductionRuleResponse(rule))
}

// ListDeductionRules handles GET /api/deduction-rules - lists all active deduction rules
func (h *DeductionHandler) ListDeductionRules(w http.ResponseWriter, r *http.Request) {
	log.Println("ListDeductionRules handler called")

	// Parse pagination parameters
	limit := 20
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	rules, total, err := h.stipendService.ListDeductionRules(limit, offset)
	if err != nil {
		log.Printf("Error fetching deduction rules: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch deduction rules: %v", err), http.StatusInternalServerError)
		return
	}

	var responses []DeductionRuleResponse
	for _, rule := range rules {
		responses = append(responses, *toDeductionRuleResponse(&rule))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rules":  responses,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// Helper functions

func toDeductionRuleResponse(rule *models.DeductionRule) *DeductionRuleResponse {
	return &DeductionRuleResponse{
		ID:                        rule.ID.String(),
		RuleName:                  rule.RuleName,
		DeductionType:             rule.DeductionType,
		Description:               rule.Description,
		BaseAmount:                rule.BaseAmount,
		MaxDeductionAmount:        rule.MaxDeductionAmount,
		MinDeductionAmount:        rule.MinDeductionAmount,
		IsApplicableToFullScholar: rule.IsApplicableToFullScholar,
		IsApplicableToSelfFunded:  rule.IsApplicableToSelfFunded,
		IsActive:                  rule.IsActive,
		AppliesMonthly:            rule.AppliesMonthly,
		AppliesAnnually:           rule.AppliesAnnually,
		IsOptional:                rule.IsOptional,
		Priority:                  rule.Priority,
		CreatedAt:                 rule.CreatedAt,
		ModifiedAt:                rule.ModifiedAt,
	}
}

func toDeductionResponse(deduction *models.Deduction) *DeductionResponse {
	var approvedBy *string
	if deduction.ApprovedBy != nil {
		s := deduction.ApprovedBy.String()
		approvedBy = &s
	}

	return &DeductionResponse{
		ID:               deduction.ID.String(),
		StudentID:        deduction.StudentID.String(),
		DeductionRuleID:  deduction.DeductionRuleID.String(),
		StipendID:        deduction.StipendID.String(),
		Amount:           deduction.Amount,
		DeductionType:    deduction.DeductionType,
		Description:      deduction.Description,
		DeductionDate:    deduction.DeductionDate,
		ProcessingStatus: deduction.ProcessingStatus,
		ApprovedBy:       approvedBy,
		ApprovalDate:     deduction.ApprovalDate,
		RejectionReason:  deduction.RejectionReason,
		CreatedAt:        deduction.CreatedAt,
		ModifiedAt:       deduction.ModifiedAt,
	}
}
