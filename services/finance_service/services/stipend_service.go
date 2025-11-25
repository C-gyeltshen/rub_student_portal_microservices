package services

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"finance_service/database"
	"finance_service/models"
)

// StipendService handles all stipend-related operations including calculation and deduction application
type StipendService struct {
	db *gorm.DB
}

// NewStipendService creates a new stipend service instance
func NewStipendService() *StipendService {
	return &StipendService{
		db: database.DB,
	}
}

// StipendCalculationResult contains the calculated stipend details
type StipendCalculationResult struct {
	BaseStipendAmount float64
	TotalDeductions   float64
	NetStipendAmount  float64
	Deductions        []DeductionDetail
}

// DeductionDetail represents a single deduction applied to a stipend
type DeductionDetail struct {
	RuleID        uuid.UUID
	RuleName      string
	DeductionType string
	Amount        float64
	Description   string
	IsOptional    bool
}

// CreateStipendForStudent creates a new stipend record for a student
func (ss *StipendService) CreateStipendForStudent(
	studentID uuid.UUID,
	stipendType string,
	amount float64,
	paymentMethod string,
	journalNumber string,
	notes string,
) (*models.Stipend, error) {
	log.Printf("Creating stipend for student %s with type %s and amount %.2f", studentID, stipendType, amount)

	if err := ss.validateStipendInput(studentID, stipendType, amount); err != nil {
		return nil, err
	}

	stipend := &models.Stipend{
		ID:            uuid.New(),
		StudentID:     studentID,
		Amount:        amount,
		StipendType:   stipendType,
		PaymentStatus: "Pending",
		PaymentMethod: paymentMethod,
		JournalNumber: journalNumber,
		Notes:         notes,
	}

	if err := ss.db.Create(stipend).Error; err != nil {
		log.Printf("Error creating stipend: %v", err)
		return nil, fmt.Errorf("failed to create stipend: %w", err)
	}

	log.Printf("Stipend created successfully with ID: %s", stipend.ID)
	return stipend, nil
}

// CalculateStipendWithDeductions calculates the net stipend amount after applying applicable deductions
func (ss *StipendService) CalculateStipendWithDeductions(
	studentID uuid.UUID,
	stipendType string,
	baseAmount float64,
) (*StipendCalculationResult, error) {
	log.Printf("Calculating stipend with deductions for student %s (type: %s, base: %.2f)", studentID, stipendType, baseAmount)

	if baseAmount <= 0 {
		return nil, fmt.Errorf("base stipend amount must be positive")
	}

	// Fetch all applicable deduction rules for this student type
	applicableRules, err := ss.getApplicableDeductionRules(stipendType)
	if err != nil {
		log.Printf("Error fetching deduction rules: %v", err)
		return nil, err
	}

	result := &StipendCalculationResult{
		BaseStipendAmount: baseAmount,
		Deductions:        []DeductionDetail{},
	}

	// Sort rules by priority (higher priority first)
	sort.Slice(applicableRules, func(i, j int) bool {
		return applicableRules[i].Priority > applicableRules[j].Priority
	})

	log.Printf("Found %d applicable deduction rules for %s students", len(applicableRules), stipendType)

	// Apply deductions in order of priority
	currentStipendAmount := baseAmount
	totalDeductions := 0.0

	for _, rule := range applicableRules {
		log.Printf("Processing deduction rule: %s (type: %s, optional: %v)", rule.RuleName, rule.DeductionType, rule.IsOptional)

		// Calculate deduction amount
		deductionAmount := ss.calculateDeductionAmount(rule, currentStipendAmount)

		if deductionAmount <= 0 {
			log.Printf("Skipping rule %s: calculated deduction is zero or negative", rule.RuleName)
			continue
		}

		// Ensure deduction doesn't exceed remaining stipend
		if deductionAmount > currentStipendAmount {
			deductionAmount = currentStipendAmount
			log.Printf("Capping deduction to remaining stipend: %.2f", deductionAmount)
		}

		// Add to result
		result.Deductions = append(result.Deductions, DeductionDetail{
			RuleID:        rule.ID,
			RuleName:      rule.RuleName,
			DeductionType: rule.DeductionType,
			Amount:        deductionAmount,
			Description:   rule.Description,
			IsOptional:    rule.IsOptional,
		})

		currentStipendAmount -= deductionAmount
		totalDeductions += deductionAmount

		log.Printf("Deduction applied: %.2f, remaining: %.2f", deductionAmount, currentStipendAmount)
	}

	result.TotalDeductions = totalDeductions
	result.NetStipendAmount = currentStipendAmount

	// Validate total deductions against stipend amount
	validationService := NewValidationService()
	deductionValidation := validationService.ValidateTotalDeductionAgainstStipend(
		totalDeductions,
		baseAmount,
		true, // allow exceed since we cap deductions to remaining stipend
	)

	if len(deductionValidation.Warnings) > 0 {
		log.Printf("[WARNING] Stipend calculation warnings: %v", deductionValidation.Warnings)
	}

	log.Printf("Stipend calculation completed: base=%.2f, deductions=%.2f, net=%.2f",
		result.BaseStipendAmount, result.TotalDeductions, result.NetStipendAmount)

	return result, nil
}

// ApplyDeductionsToStipend applies calculated deductions to a stipend record
func (ss *StipendService) ApplyDeductionsToStipend(
	stipendID uuid.UUID,
	studentID uuid.UUID,
	deductionDetails []DeductionDetail,
) ([]models.Deduction, error) {
	log.Printf("Applying %d deductions to stipend %s for student %s", len(deductionDetails), stipendID, studentID)

	var appliedDeductions []models.Deduction

	for _, detail := range deductionDetails {
		deduction := models.Deduction{
			ID:               uuid.New(),
			StudentID:        studentID,
			DeductionRuleID:  detail.RuleID,
			StipendID:        stipendID,
			Amount:           detail.Amount,
			DeductionType:    detail.DeductionType,
			Description:      detail.Description,
			ProcessingStatus: "Pending",
			DeductionDate:    time.Now(),
		}

		if err := ss.db.Create(&deduction).Error; err != nil {
			log.Printf("Error creating deduction: %v", err)
			return nil, fmt.Errorf("failed to apply deduction: %w", err)
		}

		appliedDeductions = append(appliedDeductions, deduction)
		log.Printf("Deduction created with ID: %s, amount: %.2f", deduction.ID, deduction.Amount)
	}

	return appliedDeductions, nil
}

// GetStipendByID retrieves a stipend by its ID
func (ss *StipendService) GetStipendByID(stipendID uuid.UUID) (*models.Stipend, error) {
	var stipend models.Stipend

	if err := ss.db.First(&stipend, "id = ?", stipendID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("stipend not found")
		}
		return nil, fmt.Errorf("failed to fetch stipend: %w", err)
	}

	return &stipend, nil
}

// GetStudentStipends retrieves all stipends for a student
func (ss *StipendService) GetStudentStipends(studentID uuid.UUID, limit int, offset int) ([]models.Stipend, int64, error) {
	var stipends []models.Stipend
	var total int64

	if err := ss.db.Where("student_id = ?", studentID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count stipends: %w", err)
	}

	if err := ss.db.
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&stipends).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch stipends: %w", err)
	}

	return stipends, total, nil
}

// UpdateStipendPaymentStatus updates the payment status of a stipend
func (ss *StipendService) UpdateStipendPaymentStatus(
	stipendID uuid.UUID,
	status string,
	paymentDate *time.Time,
) error {
	log.Printf("Updating stipend %s payment status to %s", stipendID, status)

	// Validate status
	validStatuses := map[string]bool{"Pending": true, "Processed": true, "Failed": true}
	if !validStatuses[status] {
		return fmt.Errorf("invalid payment status: %s", status)
	}

	update := map[string]interface{}{
		"payment_status": status,
	}

	if paymentDate != nil {
		update["payment_date"] = paymentDate
	}

	if err := ss.db.Model(&models.Stipend{}).Where("id = ?", stipendID).Updates(update).Error; err != nil {
		log.Printf("Error updating stipend status: %v", err)
		return fmt.Errorf("failed to update stipend status: %w", err)
	}

	return nil
}

// GetStipendDeductions retrieves all deductions for a stipend
func (ss *StipendService) GetStipendDeductions(stipendID uuid.UUID) ([]models.Deduction, error) {
	var deductions []models.Deduction

	if err := ss.db.
		Where("stipend_id = ?", stipendID).
		Order("deduction_date DESC").
		Find(&deductions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch deductions: %w", err)
	}

	return deductions, nil
}

// CalculateMonthlyStipendForStudent calculates the monthly stipend for a student
// This is typically 1/12th of the annual stipend, minus applicable monthly deductions
func (ss *StipendService) CalculateMonthlyStipendForStudent(
	studentID uuid.UUID,
	stipendType string,
	annualAmount float64,
) (*StipendCalculationResult, error) {
	log.Printf("Calculating monthly stipend for student %s (annual: %.2f)", studentID, annualAmount)

	monthlyAmount := annualAmount / 12
	log.Printf("Monthly base amount: %.2f", monthlyAmount)

	return ss.CalculateStipendWithDeductions(studentID, stipendType, monthlyAmount)
}

// CalculateAnnualStipendForStudent calculates the annual stipend for a student
// This includes all applicable annual deductions
func (ss *StipendService) CalculateAnnualStipendForStudent(
	studentID uuid.UUID,
	stipendType string,
	annualAmount float64,
) (*StipendCalculationResult, error) {
	log.Printf("Calculating annual stipend for student %s (amount: %.2f)", studentID, annualAmount)

	return ss.CalculateStipendWithDeductions(studentID, stipendType, annualAmount)
}

// Helper functions

// validateStipendInput validates the input for stipend creation
func (ss *StipendService) validateStipendInput(studentID uuid.UUID, stipendType string, amount float64) error {
	validationService := NewValidationService()

	// Validate stipend input using validation service
	validationResult := validationService.ValidateStipendInput(
		studentID,
		stipendType,
		amount,
		"", // journal number is validated separately in CreateStipendForStudent
	)

	if !validationResult.IsValid {
		validationService.LogValidationResult("ValidateStipendInput", validationResult)
		return fmt.Errorf("%s", validationService.FormatValidationError(validationResult))
	}

	// Log warnings if any
	if len(validationResult.Warnings) > 0 {
		log.Printf("[WARNING] %s", validationService.FormatValidationWarnings(validationResult))
	}

	return nil
}

// getApplicableDeductionRules fetches deduction rules applicable to a student type
func (ss *StipendService) getApplicableDeductionRules(stipendType string) ([]models.DeductionRule, error) {
	var rules []models.DeductionRule

	query := ss.db.Where("is_active = ?", true)

	if stipendType == "full-scholarship" {
		query = query.Where("is_applicable_to_full_scholar = ?", true)
	} else if stipendType == "self-funded" {
		query = query.Where("is_applicable_to_self_funded = ?", true)
	} else {
		return nil, fmt.Errorf("invalid stipend type: %s", stipendType)
	}

	if err := query.Order("priority DESC").Find(&rules).Error; err != nil {
		log.Printf("Error fetching deduction rules: %v", err)
		return nil, fmt.Errorf("failed to fetch deduction rules: %w", err)
	}

	return rules, nil
}

// calculateDeductionAmount calculates the deduction amount based on the rule
func (ss *StipendService) calculateDeductionAmount(rule models.DeductionRule, currentStipend float64) float64 {
	// Start with base amount
	amount := rule.BaseAmount

	// Ensure it doesn't exceed max deduction
	if amount > rule.MaxDeductionAmount {
		amount = rule.MaxDeductionAmount
	}

	// Ensure it meets minimum deduction
	if amount < rule.MinDeductionAmount {
		amount = rule.MinDeductionAmount
	}

	// Ensure it doesn't exceed current stipend
	if amount > currentStipend {
		amount = currentStipend
	}

	// If the rule is optional, it might be skipped (handled at a higher level)
	// This function just calculates the amount

	return amount
}

// GetDeductionRuleByID retrieves a deduction rule by ID
func (ss *StipendService) GetDeductionRuleByID(ruleID uuid.UUID) (*models.DeductionRule, error) {
	var rule models.DeductionRule

	if err := ss.db.First(&rule, "id = ?", ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("deduction rule not found")
		}
		return nil, fmt.Errorf("failed to fetch deduction rule: %w", err)
	}

	return &rule, nil
}

// ListDeductionRules retrieves all active deduction rules
func (ss *StipendService) ListDeductionRules(limit int, offset int) ([]models.DeductionRule, int64, error) {
	var rules []models.DeductionRule
	var total int64

	if err := ss.db.Model(&models.DeductionRule{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rules: %w", err)
	}

	if err := ss.db.
		Where("is_active = ?", true).
		Order("priority DESC, rule_name ASC").
		Limit(limit).
		Offset(offset).
		Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch rules: %w", err)
	}

	return rules, total, nil
}

// CreateDeductionRule creates a new deduction rule
func (ss *StipendService) CreateDeductionRule(rule *models.DeductionRule) error {
	log.Printf("Creating deduction rule: %s", rule.RuleName)

	if rule.RuleName == "" {
		return fmt.Errorf("rule name is required")
	}

	if rule.DeductionType == "" {
		return fmt.Errorf("deduction type is required")
	}

	if rule.BaseAmount < 0 {
		return fmt.Errorf("base amount cannot be negative")
	}

	if rule.MaxDeductionAmount < 0 {
		return fmt.Errorf("max deduction amount cannot be negative")
	}

	rule.ID = uuid.New()

	if err := ss.db.Create(rule).Error; err != nil {
		log.Printf("Error creating deduction rule: %v", err)
		return fmt.Errorf("failed to create deduction rule: %w", err)
	}

	log.Printf("Deduction rule created with ID: %s", rule.ID)
	return nil
}

// UpdateDeductionRule updates an existing deduction rule
func (ss *StipendService) UpdateDeductionRule(ruleID uuid.UUID, updates map[string]interface{}) (*models.DeductionRule, error) {
	log.Printf("Updating deduction rule: %s", ruleID)

	// Fetch the existing rule first
	var rule models.DeductionRule
	if err := ss.db.First(&rule, "id = ?", ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("deduction rule not found")
		}
		return nil, fmt.Errorf("failed to fetch deduction rule: %w", err)
	}

	// Validate updates
	if ruleName, ok := updates["rule_name"]; ok {
		if ruleName == "" {
			return nil, fmt.Errorf("rule name cannot be empty")
		}
		// Check for duplicate rule name (excluding current rule)
		var existingRule models.DeductionRule
		if err := ss.db.Where("rule_name = ? AND id != ?", ruleName, ruleID).First(&existingRule).Error; err == nil {
			return nil, fmt.Errorf("rule name already exists")
		}
	}

	if deductionType, ok := updates["deduction_type"]; ok && deductionType == "" {
		return nil, fmt.Errorf("deduction type cannot be empty")
	}

	if baseAmount, ok := updates["base_amount"]; ok {
		if ba, err := toFloat64(baseAmount); err != nil || ba < 0 {
			return nil, fmt.Errorf("base amount must be non-negative")
		}
	}

	if maxAmount, ok := updates["max_deduction_amount"]; ok {
		if ma, err := toFloat64(maxAmount); err != nil || ma < 0 {
			return nil, fmt.Errorf("max deduction amount must be non-negative")
		}
	}

	if minAmount, ok := updates["min_deduction_amount"]; ok {
		if mi, err := toFloat64(minAmount); err != nil || mi < 0 {
			return nil, fmt.Errorf("min deduction amount must be non-negative")
		}
	}

	// Perform the update
	if err := ss.db.Model(&rule).Updates(updates).Error; err != nil {
		log.Printf("Error updating deduction rule: %v", err)
		return nil, fmt.Errorf("failed to update deduction rule: %w", err)
	}

	log.Printf("Deduction rule %s updated successfully", ruleID)
	return &rule, nil
}

// DeleteDeductionRule soft-deletes a deduction rule by marking it as inactive
// Note: We soft-delete to maintain referential integrity with deductions
func (ss *StipendService) DeleteDeductionRule(ruleID uuid.UUID) error {
	log.Printf("Deleting deduction rule: %s", ruleID)

	var rule models.DeductionRule
	if err := ss.db.First(&rule, "id = ?", ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("deduction rule not found")
		}
		return fmt.Errorf("failed to fetch deduction rule: %w", err)
	}

	// Soft delete by marking as inactive
	if err := ss.db.Model(&rule).Update("is_active", false).Error; err != nil {
		log.Printf("Error deleting deduction rule: %v", err)
		return fmt.Errorf("failed to delete deduction rule: %w", err)
	}

	log.Printf("Deduction rule %s deleted successfully", ruleID)
	return nil
}

// Helper function to convert interface{} to float64
func toFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		var f float64
		if _, err := fmt.Sscanf(v, "%f", &f); err != nil {
			return 0, err
		}
		return f, nil
	default:
		return 0, fmt.Errorf("unsupported type for conversion to float64")
	}
}

// GetStudentStipendsWithPagination retrieves all stipends for a student with pagination (returns converted service types)
func (ss *StipendService) GetStudentStipendsWithPagination(studentID uuid.UUID, limit int, offset int) ([]*Stipend, int64, error) {
	var modelStipends []models.Stipend
	var total int64

	if err := ss.db.Where("student_id = ?", studentID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count stipends: %w", err)
	}

	if err := ss.db.
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&modelStipends).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch stipends: %w", err)
	}

	stipends := make([]*Stipend, len(modelStipends))
	for i, ms := range modelStipends {
		stipends[i] = convertModelStipendToService(&ms)
	}

	return stipends, total, nil
}

// UpdateStipendPaymentStatus updates the payment status of a stipend and returns the updated stipend
func (ss *StipendService) UpdateStipendPaymentStatusWithReturn(
	stipendID uuid.UUID,
	status string,
	paymentDate *time.Time,
) (*Stipend, error) {
	log.Printf("Updating stipend %s payment status to %s", stipendID, status)

	// Validate status
	validStatuses := map[string]bool{"Pending": true, "Processed": true, "Failed": true}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid payment status: %s", status)
	}

	update := map[string]interface{}{
		"payment_status": status,
	}

	if paymentDate != nil {
		update["payment_date"] = paymentDate
	}

	if err := ss.db.Model(&models.Stipend{}).Where("id = ?", stipendID).Updates(update).Error; err != nil {
		log.Printf("Error updating stipend status: %v", err)
		return nil, fmt.Errorf("failed to update stipend status: %w", err)
	}

	// Retrieve and return the updated stipend
	var stipend models.Stipend
	if err := ss.db.First(&stipend, "id = ?", stipendID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated stipend: %w", err)
	}

	return convertModelStipendToService(&stipend), nil
}

