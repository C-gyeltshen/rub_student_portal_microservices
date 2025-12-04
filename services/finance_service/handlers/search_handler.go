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
)

// SearchHandler handles search and filter requests
type SearchHandler struct {
	searchService *services.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler() *SearchHandler {
	return &SearchHandler{
		searchService: services.NewSearchService(database.DB),
	}
}

// SearchStipends handles searching stipends with filters and pagination
func (sh *SearchHandler) SearchStipends(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	studentID := r.URL.Query().Get("student_id")
	status := r.URL.Query().Get("status")
	stipendType := r.URL.Query().Get("stipend_type")
	paymentStatus := r.URL.Query().Get("payment_status")
	startDateStr := r.URL.Query().Get("start_date") // RFC3339 format
	endDateStr := r.URL.Query().Get("end_date")     // RFC3339 format
	minAmountStr := r.URL.Query().Get("min_amount")
	maxAmountStr := r.URL.Query().Get("max_amount")
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

	// Parse amounts
	minAmount := 0.0
	if m, err := strconv.ParseFloat(minAmountStr, 64); err == nil && m > 0 {
		minAmount = m
	}
	maxAmount := 0.0
	if m, err := strconv.ParseFloat(maxAmountStr, 64); err == nil && m > 0 {
		maxAmount = m
	}

	// Parse dates
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	params := services.SearchStipendParams{
		StudentID:     studentID,
		Status:        status,
		StipendType:   stipendType,
		PaymentStatus: paymentStatus,
		StartDate:     startDate,
		EndDate:       endDate,
		MinAmount:     minAmount,
		MaxAmount:     maxAmount,
		Limit:         limit,
		Offset:        offset,
	}

	stipends, total, err := sh.searchService.SearchStipends(params)
	if err != nil {
		log.Printf("Error searching stipends: %v", err)
		http.Error(w, fmt.Sprintf("Failed to search stipends: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       stipends,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"count":      len(stipends),
	})
}

// SearchDeductionRules handles searching deduction rules with filters and pagination
func (sh *SearchHandler) SearchDeductionRules(w http.ResponseWriter, r *http.Request) {
	ruleName := r.URL.Query().Get("rule_name")
	deductionType := r.URL.Query().Get("deduction_type")
	isActiveStr := r.URL.Query().Get("is_active")
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

	// Parse is_active filter
	var isActive *bool
	if isActiveStr != "" {
		active := isActiveStr == "true"
		isActive = &active
	}

	params := services.SearchDeductionRulesParams{
		RuleName:      ruleName,
		DeductionType: deductionType,
		IsActive:      isActive,
		Limit:         limit,
		Offset:        offset,
	}

	rules, total, err := sh.searchService.SearchDeductionRules(params)
	if err != nil {
		log.Printf("Error searching deduction rules: %v", err)
		http.Error(w, fmt.Sprintf("Failed to search deduction rules: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       rules,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"count":      len(rules),
	})
}

// SearchTransactions handles searching transactions with filters and pagination
func (sh *SearchHandler) SearchTransactions(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("student_id")
	stipendID := r.URL.Query().Get("stipend_id")
	status := r.URL.Query().Get("status")
	transactionType := r.URL.Query().Get("transaction_type")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	minAmountStr := r.URL.Query().Get("min_amount")
	maxAmountStr := r.URL.Query().Get("max_amount")
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

	// Parse amounts
	minAmount := 0.0
	if m, err := strconv.ParseFloat(minAmountStr, 64); err == nil && m > 0 {
		minAmount = m
	}
	maxAmount := 0.0
	if m, err := strconv.ParseFloat(maxAmountStr, 64); err == nil && m > 0 {
		maxAmount = m
	}

	// Parse dates
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	params := services.SearchTransactionsParams{
		StudentID:       studentID,
		StipendID:       stipendID,
		Status:          status,
		TransactionType: transactionType,
		StartDate:       startDate,
		EndDate:         endDate,
		MinAmount:       minAmount,
		MaxAmount:       maxAmount,
		Limit:           limit,
		Offset:          offset,
	}

	transactions, total, err := sh.searchService.SearchTransactions(params)
	if err != nil {
		log.Printf("Error searching transactions: %v", err)
		http.Error(w, fmt.Sprintf("Failed to search transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       transactions,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"count":      len(transactions),
	})
}
