package services

import (
	"finance_service/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// SearchService handles searching and filtering operations
type SearchService struct {
	db *gorm.DB
}

// NewSearchService creates a new search service
func NewSearchService(db *gorm.DB) *SearchService {
	return &SearchService{db: db}
}

// SearchStipendParams holds search parameters for stipends
type SearchStipendParams struct {
	StudentID      string
	Status         string
	StipendType    string
	StartDate      *time.Time
	EndDate        *time.Time
	MinAmount      float64
	MaxAmount      float64
	PaymentStatus  string
	Limit          int
	Offset         int
}

// SearchStipends searches stipends with multiple filters and pagination
func (ss *SearchService) SearchStipends(params SearchStipendParams) ([]models.Stipend, int64, error) {
	var stipends []models.Stipend
	var total int64

	query := ss.db

	// Apply filters
	if params.StudentID != "" {
		query = query.Where("student_id = ?", params.StudentID)
	}
	if params.Status != "" {
		query = query.Where("payment_status = ?", params.Status)
	}
	if params.StipendType != "" {
		query = query.Where("stipend_type = ?", params.StipendType)
	}
	if params.PaymentStatus != "" {
		query = query.Where("payment_status = ?", params.PaymentStatus)
	}
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", params.EndDate)
	}
	if params.MinAmount > 0 {
		query = query.Where("amount >= ?", params.MinAmount)
	}
	if params.MaxAmount > 0 {
		query = query.Where("amount <= ?", params.MaxAmount)
	}

	// Get total count
	if err := query.Model(&models.Stipend{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count stipends: %w", err)
	}

	// Set default limit
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100 // Max 100 per page
	}

	// Get paginated results
	if err := query.Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&stipends).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search stipends: %w", err)
	}

	return stipends, total, nil
}

// SearchDeductionRulesParams holds search parameters for deduction rules
type SearchDeductionRulesParams struct {
	RuleName      string
	DeductionType string
	IsActive      *bool
	Limit         int
	Offset        int
}

// SearchDeductionRules searches deduction rules with filters and pagination
func (ss *SearchService) SearchDeductionRules(params SearchDeductionRulesParams) ([]models.DeductionRule, int64, error) {
	var rules []models.DeductionRule
	var total int64

	query := ss.db

	// Apply filters
	if params.RuleName != "" {
		query = query.Where("rule_name ILIKE ?", "%"+params.RuleName+"%")
	}
	if params.DeductionType != "" {
		query = query.Where("deduction_type = ?", params.DeductionType)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Get total count
	if err := query.Model(&models.DeductionRule{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count deduction rules: %w", err)
	}

	// Set default limit
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Get paginated results
	if err := query.Order("priority ASC, created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search deduction rules: %w", err)
	}

	return rules, total, nil
}

// SearchTransactionsParams holds search parameters for transactions
type SearchTransactionsParams struct {
	StudentID      string
	StipendID      string
	Status         string
	TransactionType string
	StartDate      *time.Time
	EndDate        *time.Time
	MinAmount      float64
	MaxAmount      float64
	Limit          int
	Offset         int
}

// SearchTransactions searches transactions with filters and pagination
func (ss *SearchService) SearchTransactions(params SearchTransactionsParams) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := ss.db

	// Apply filters
	if params.StudentID != "" {
		query = query.Where("student_id = ?", params.StudentID)
	}
	if params.StipendID != "" {
		query = query.Where("stipend_id = ?", params.StipendID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.TransactionType != "" {
		query = query.Where("transaction_type = ?", params.TransactionType)
	}
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", params.EndDate)
	}
	if params.MinAmount > 0 {
		query = query.Where("amount >= ?", params.MinAmount)
	}
	if params.MaxAmount > 0 {
		query = query.Where("amount <= ?", params.MaxAmount)
	}

	// Get total count
	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Set default limit
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Get paginated results
	if err := query.Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search transactions: %w", err)
	}

	return transactions, total, nil
}
