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

// StipendHandler handles HTTP requests for stipend operations
type StipendHandler struct {
	stipendService *services.StipendService
}

// NewStipendHandler creates a new stipend handler
func NewStipendHandler() *StipendHandler {
	return &StipendHandler{
		stipendService: services.NewStipendService(),
	}
}

// CreateStipendRequest represents the request body for creating a stipend
type CreateStipendRequest struct {
	StudentID     string  `json:"student_id"`
	StipendType   string  `json:"stipend_type"` // full-scholarship or self-funded
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	JournalNumber string  `json:"journal_number"`
	Notes         string  `json:"notes,omitempty"`
}

// CreateStipendWithDeductionsRequest represents the request for creating a stipend with automatic deduction calculation
type CreateStipendWithDeductionsRequest struct {
	StudentID     string  `json:"student_id"`
	StipendType   string  `json:"stipend_type"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	JournalNumber string  `json:"journal_number"`
	Notes         string  `json:"notes,omitempty"`
}

// StipendResponse represents the response for a stipend
type StipendResponse struct {
	ID            string     `json:"id"`
	StudentID     string     `json:"student_id"`
	Amount        float64    `json:"amount"`
	StipendType   string     `json:"stipend_type"`
	PaymentDate   *time.Time `json:"payment_date"`
	PaymentStatus string     `json:"payment_status"`
	PaymentMethod string     `json:"payment_method"`
	JournalNumber string     `json:"journal_number"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
	ModifiedAt    time.Time  `json:"modified_at"`
}

// CalculationResponse represents the response for stipend calculation with deductions
type CalculationResponse struct {
	BaseStipendAmount float64                   `json:"base_stipend_amount"`
	TotalDeductions   float64                   `json:"total_deductions"`
	NetStipendAmount  float64                   `json:"net_stipend_amount"`
	Deductions        []DeductionDetailResponse `json:"deductions"`
}

// DeductionDetailResponse represents a deduction detail in the response
type DeductionDetailResponse struct {
	RuleID        string  `json:"rule_id"`
	RuleName      string  `json:"rule_name"`
	DeductionType string  `json:"deduction_type"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	IsOptional    bool    `json:"is_optional"`
}

// CreateStipend handles POST /api/stipends - creates a new stipend
func (h *StipendHandler) CreateStipend(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateStipend handler called")

	var req CreateStipendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Parse student ID
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}

	// Create stipend
	stipend, err := h.stipendService.CreateStipendForStudent(
		studentID,
		req.StipendType,
		req.Amount,
		req.PaymentMethod,
		req.JournalNumber,
		req.Notes,
	)
	if err != nil {
		log.Printf("Error creating stipend: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create stipend: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toStipendResponse(stipend))
}

// CalculateStipendWithDeductions handles POST /api/stipends/calculate - calculates stipend with deductions
func (h *StipendHandler) CalculateStipendWithDeductions(w http.ResponseWriter, r *http.Request) {
	log.Println("CalculateStipendWithDeductions handler called")

	var req CreateStipendWithDeductionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Parse student ID
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}

	// Calculate stipend with deductions
	result, err := h.stipendService.CalculateStipendWithDeductions(
		studentID,
		req.StipendType,
		req.Amount,
	)
	if err != nil {
		log.Printf("Error calculating stipend: %v", err)
		http.Error(w, fmt.Sprintf("Failed to calculate stipend: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toCalculationResponse(result))
}

// GetStipend handles GET /api/stipends/{stipendID} - retrieves a stipend by ID
func (h *StipendHandler) GetStipend(w http.ResponseWriter, r *http.Request) {
	log.Println("GetStipend handler called")

	stipendIDStr := chi.URLParam(r, "stipendID")
	stipendID, err := uuid.Parse(stipendIDStr)
	if err != nil {
		http.Error(w, "Invalid stipend ID format", http.StatusBadRequest)
		return
	}

	stipend, err := h.stipendService.GetStipendByID(stipendID)
	if err != nil {
		log.Printf("Error fetching stipend: %v", err)
		http.Error(w, "Stipend not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toStipendResponse(stipend))
}

// GetStudentStipends handles GET /api/students/{studentID}/stipends - retrieves stipends for a student
func (h *StipendHandler) GetStudentStipends(w http.ResponseWriter, r *http.Request) {
	log.Println("GetStudentStipends handler called")

	studentIDStr := chi.URLParam(r, "studentID")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	limit := 10
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

	stipends, total, err := h.stipendService.GetStudentStipends(studentID, limit, offset)
	if err != nil {
		log.Printf("Error fetching stipends: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch stipends: %v", err), http.StatusInternalServerError)
		return
	}

	var responses []StipendResponse
	for _, stipend := range stipends {
		responses = append(responses, *toStipendResponse(&stipend))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"stipends": responses,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// UpdateStipendPaymentStatus handles PATCH /api/stipends/{stipendID}/payment-status - updates payment status
func (h *StipendHandler) UpdateStipendPaymentStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateStipendPaymentStatus handler called")

	stipendIDStr := chi.URLParam(r, "stipendID")
	stipendID, err := uuid.Parse(stipendIDStr)
	if err != nil {
		http.Error(w, "Invalid stipend ID format", http.StatusBadRequest)
		return
	}

	var req struct {
		Status      string `json:"status"`
		PaymentDate string `json:"payment_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	var paymentDate *time.Time
	if req.PaymentDate != "" {
		t, err := time.Parse(time.RFC3339, req.PaymentDate)
		if err != nil {
			http.Error(w, "Invalid payment date format", http.StatusBadRequest)
			return
		}
		paymentDate = &t
	}

	if err := h.stipendService.UpdateStipendPaymentStatus(stipendID, req.Status, paymentDate); err != nil {
		log.Printf("Error updating stipend status: %v", err)
		http.Error(w, fmt.Sprintf("Failed to update stipend: %v", err), http.StatusBadRequest)
		return
	}

	// Fetch updated stipend
	stipend, err := h.stipendService.GetStipendByID(stipendID)
	if err != nil {
		log.Printf("Error fetching updated stipend: %v", err)
		http.Error(w, "Failed to fetch updated stipend", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toStipendResponse(stipend))
}

// GetStipendDeductions handles GET /api/stipends/{stipendID}/deductions - retrieves deductions for a stipend
func (h *StipendHandler) GetStipendDeductions(w http.ResponseWriter, r *http.Request) {
	log.Println("GetStipendDeductions handler called")

	stipendIDStr := chi.URLParam(r, "stipendID")
	stipendID, err := uuid.Parse(stipendIDStr)
	if err != nil {
		http.Error(w, "Invalid stipend ID format", http.StatusBadRequest)
		return
	}

	deductions, err := h.stipendService.GetStipendDeductions(stipendID)
	if err != nil {
		log.Printf("Error fetching deductions: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch deductions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deductions": deductions,
		"count":      len(deductions),
	})
}

// CalculateMonthlyStipend handles POST /api/stipends/calculate/monthly - calculates monthly stipend
func (h *StipendHandler) CalculateMonthlyStipend(w http.ResponseWriter, r *http.Request) {
	log.Println("CalculateMonthlyStipend handler called")

	var req struct {
		StudentID    string  `json:"student_id"`
		StipendType  string  `json:"stipend_type"`
		AnnualAmount float64 `json:"annual_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}

	result, err := h.stipendService.CalculateMonthlyStipendForStudent(
		studentID,
		req.StipendType,
		req.AnnualAmount,
	)
	if err != nil {
		log.Printf("Error calculating monthly stipend: %v", err)
		http.Error(w, fmt.Sprintf("Failed to calculate monthly stipend: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toCalculationResponse(result))
}

// CalculateAnnualStipend handles POST /api/stipends/calculate/annual - calculates annual stipend
func (h *StipendHandler) CalculateAnnualStipend(w http.ResponseWriter, r *http.Request) {
	log.Println("CalculateAnnualStipend handler called")

	var req struct {
		StudentID    string  `json:"student_id"`
		StipendType  string  `json:"stipend_type"`
		AnnualAmount float64 `json:"annual_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}

	result, err := h.stipendService.CalculateAnnualStipendForStudent(
		studentID,
		req.StipendType,
		req.AnnualAmount,
	)
	if err != nil {
		log.Printf("Error calculating annual stipend: %v", err)
		http.Error(w, fmt.Sprintf("Failed to calculate annual stipend: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toCalculationResponse(result))
}

// Helper functions

func toStipendResponse(stipend *models.Stipend) *StipendResponse {
	return &StipendResponse{
		ID:            stipend.ID.String(),
		StudentID:     stipend.StudentID.String(),
		Amount:        stipend.Amount,
		StipendType:   stipend.StipendType,
		PaymentDate:   stipend.PaymentDate,
		PaymentStatus: stipend.PaymentStatus,
		PaymentMethod: stipend.PaymentMethod,
		JournalNumber: stipend.JournalNumber,
		Notes:         stipend.Notes,
		CreatedAt:     stipend.CreatedAt,
		ModifiedAt:    stipend.ModifiedAt,
	}
}

func toCalculationResponse(result *services.StipendCalculationResult) *CalculationResponse {
	var deductions []DeductionDetailResponse
	for _, d := range result.Deductions {
		deductions = append(deductions, DeductionDetailResponse{
			RuleID:        d.RuleID.String(),
			RuleName:      d.RuleName,
			DeductionType: d.DeductionType,
			Amount:        d.Amount,
			Description:   d.Description,
			IsOptional:    d.IsOptional,
		})
	}

	return &CalculationResponse{
		BaseStipendAmount: result.BaseStipendAmount,
		TotalDeductions:   result.TotalDeductions,
		NetStipendAmount:  result.NetStipendAmount,
		Deductions:        deductions,
	}
}
