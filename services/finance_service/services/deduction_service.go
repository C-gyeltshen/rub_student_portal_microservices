package services

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"finance_service/database"
	"finance_service/models"
)

// DeductionService handles all deduction-related operations
type DeductionService struct {
	db *gorm.DB
}

// NewDeductionService creates a new deduction service instance
func NewDeductionService() *DeductionService {
	return &DeductionService{
		db: database.DB,
	}
}

// CreateDeductionRule creates a new deduction rule
func (ds *DeductionService) CreateDeductionRule(
	ruleName string,
	deductionType string,
	defaultAmount float64,
	minAmount float64,
	maxAmount float64,
	frequency string,
	isMandatory bool,
	applicableTo string,
	priority int,
	description string,
) (*DeductionRule, error) {
	log.Printf("Creating deduction rule: %s (type: %s)", ruleName, deductionType)

	// Validate input using validation service
	validationService := NewValidationService()
	validationResult := validationService.ValidateDeductionRuleInput(
		ruleName, deductionType, defaultAmount, minAmount, maxAmount,
	)

	if !validationResult.IsValid {
		validationService.LogValidationResult("CreateDeductionRule", validationResult)
		return nil, fmt.Errorf("%s", validationService.FormatValidationError(validationResult))
	}

	// Log warnings if any
	if len(validationResult.Warnings) > 0 {
		log.Printf("[WARNING] %s", validationService.FormatValidationWarnings(validationResult))
	}

	// Create the model
	modelRule := &models.DeductionRule{
		ID:                       uuid.New(),
		RuleName:                 ruleName,
		DeductionType:            deductionType,
		BaseAmount:               defaultAmount,
		MinDeductionAmount:       minAmount,
		MaxDeductionAmount:       maxAmount,
		IsActive:                 true,
		Priority:                 priority,
		Description:              description,
	}

	// Set applicable types based on the applicableTo parameter
	if applicableTo == "All" || applicableTo == "FullScholarship" {
		modelRule.IsApplicableToFullScholar = true
	}
	if applicableTo == "All" || applicableTo == "SelfFunded" {
		modelRule.IsApplicableToSelfFunded = true
	}

	// Set frequency-based flags
	if frequency == "Monthly" {
		modelRule.AppliesMonthly = true
	} else if frequency == "Annual" {
		modelRule.AppliesAnnually = true
	}

	// Set mandatory flag
	modelRule.IsOptional = !isMandatory


	if err := ds.db.Create(modelRule).Error; err != nil {
		log.Printf("Error creating deduction rule: %v", err)
		return nil, fmt.Errorf("failed to create deduction rule: %w", err)
	}

	log.Printf("Deduction rule created with ID: %s", modelRule.ID)
	return convertModelDeductionRuleToService(modelRule), nil
}

// GetDeductionRuleByID retrieves a deduction rule by ID
func (ds *DeductionService) GetDeductionRuleByID(ruleID uuid.UUID) (*DeductionRule, error) {
	var rule models.DeductionRule

	if err := ds.db.First(&rule, "id = ?", ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("deduction rule not found")
		}
		return nil, fmt.Errorf("failed to fetch deduction rule: %w", err)
	}

	return convertModelDeductionRuleToService(&rule), nil
}

// ListDeductionRulesWithPagination lists all deduction rules with optional filters and pagination
func (ds *DeductionService) ListDeductionRulesWithPagination(applicableTo string, isActive bool, limit int, offset int) ([]*DeductionRule, int64, error) {
	var modelRules []models.DeductionRule
	var total int64

	query := ds.db
	if isActive {
		query = query.Where("is_active = ?", true)
	}
	if applicableTo != "" {
		query = query.Where("applicable_to = ?", applicableTo)
	}

	// Count with explicit table model
	if err := query.Model(&models.DeductionRule{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rules: %w", err)
	}

	if err := query.
		Order("priority DESC, rule_name ASC").
		Limit(limit).
		Offset(offset).
		Find(&modelRules).Error; err != nil {
		log.Printf("Error listing deduction rules: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch rules: %w", err)
	}

	rules := make([]*DeductionRule, len(modelRules))
	for i, mr := range modelRules {
		rules[i] = convertModelDeductionRuleToService(&mr)
	}

	return rules, total, nil
}

// ApplyDeductions applies deductions to a stipend
func (ds *DeductionService) ApplyDeductions(
	studentID uuid.UUID,
	stipendID uuid.UUID,
	stipendType string,
	baseAmount float64,
	ruleIDs []uuid.UUID,
) ([]*Deduction, float64, error) {
	log.Printf("Applying deductions to stipend %s for student %s", stipendID, studentID)

	var appliedDeductions []*Deduction
	totalAmount := 0.0

	// If specific rules are provided, use them; otherwise get all applicable rules
	var rules []models.DeductionRule

	if len(ruleIDs) > 0 {
		// Fetch specific rules
		if err := ds.db.Where("id IN ?", ruleIDs).Order("priority DESC").Find(&rules).Error; err != nil {
			return nil, 0, fmt.Errorf("failed to fetch deduction rules: %w", err)
		}
	} else {
		// Get all applicable rules for this student type
		query := ds.db.Where("is_active = ?", true)

		if stipendType == "full-scholarship" {
			query = query.Where("is_applicable_to_full_scholar = ?", true)
		} else if stipendType == "self-funded" {
			query = query.Where("is_applicable_to_self_funded = ?", true)
		}

		if err := query.Order("priority DESC").Find(&rules).Error; err != nil {
			return nil, 0, fmt.Errorf("failed to fetch deduction rules: %w", err)
		}
	}

	// Apply each deduction
	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		// Calculate deduction amount
		deductionAmount := rule.BaseAmount
		if deductionAmount > rule.MaxDeductionAmount {
			deductionAmount = rule.MaxDeductionAmount
		}
		if deductionAmount < rule.MinDeductionAmount {
			deductionAmount = rule.MinDeductionAmount
		}

		// Ensure deduction doesn't exceed base amount
		if deductionAmount > baseAmount {
			deductionAmount = baseAmount
		}

		if deductionAmount <= 0 {
			continue
		}

		// Create deduction record
		modelDeduction := &models.Deduction{
			ID:               uuid.New(),
			StudentID:        studentID,
			DeductionRuleID:  rule.ID,
			StipendID:        stipendID,
			Amount:           deductionAmount,
			DeductionType:    rule.DeductionType,
			Description:      rule.Description,
			ProcessingStatus: "Pending",
			DeductionDate:    time.Now(),
		}

		if err := ds.db.Create(modelDeduction).Error; err != nil {
			log.Printf("Error creating deduction: %v", err)
			return nil, 0, fmt.Errorf("failed to apply deduction: %w", err)
		}

		appliedDeductions = append(appliedDeductions, convertModelDeductionToService(modelDeduction))
		totalAmount += deductionAmount

		log.Printf("Deduction applied: %s, amount: %.2f", rule.RuleName, deductionAmount)
	}

	return appliedDeductions, totalAmount, nil
}

// CreateDeduction creates a new deduction record
func (ds *DeductionService) CreateDeduction(
	studentID uuid.UUID,
	ruleID uuid.UUID,
	stipendID uuid.UUID,
	amount float64,
	deductionType string,
	description string,
) (*Deduction, error) {
	log.Printf("Creating deduction for student %s", studentID)

	// Validate using validation service
	validationService := NewValidationService()

	// Validate the deduction amount against the rule
	validationResult := validationService.ValidateDeductionAmount(amount, ruleID, "Deduction amount")
	if !validationResult.IsValid {
		validationService.LogValidationResult("CreateDeduction", validationResult)
		return nil, fmt.Errorf("%s", validationService.FormatValidationError(validationResult))
	}

	// Log warnings if any
	if len(validationResult.Warnings) > 0 {
		log.Printf("[WARNING] %s", validationService.FormatValidationWarnings(validationResult))
	}

	modelDeduction := &models.Deduction{
		ID:               uuid.New(),
		StudentID:        studentID,
		DeductionRuleID:  ruleID,
		StipendID:        stipendID,
		Amount:           amount,
		DeductionType:    deductionType,
		Description:      description,
		ProcessingStatus: "Pending",
		DeductionDate:    time.Now(),
	}

	if err := ds.db.Create(modelDeduction).Error; err != nil {
		log.Printf("Error creating deduction: %v", err)
		return nil, fmt.Errorf("failed to create deduction: %w", err)
	}

	log.Printf("Deduction created with ID: %s", modelDeduction.ID)
	return convertModelDeductionToService(modelDeduction), nil
}

// GetDeductionByID retrieves a deduction by ID
func (ds *DeductionService) GetDeductionByID(deductionID uuid.UUID) (*Deduction, error) {
	var deduction models.Deduction

	if err := ds.db.First(&deduction, "id = ?", deductionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("deduction not found")
		}
		return nil, fmt.Errorf("failed to fetch deduction: %w", err)
	}

	return convertModelDeductionToService(&deduction), nil
}

// GetStipendDeductionsWithPagination retrieves all deductions for a stipend with pagination
func (ds *DeductionService) GetStipendDeductionsWithPagination(stipendID uuid.UUID, limit int, offset int) ([]*Deduction, float64, error) {
	var modelDeductions []models.Deduction
	var total int64

	if err := ds.db.Model(&models.Deduction{}).Where("stipend_id = ?", stipendID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count deductions: %w", err)
	}

	if err := ds.db.
		Where("stipend_id = ?", stipendID).
		Order("deduction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&modelDeductions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch deductions: %w", err)
	}

	deductions := make([]*Deduction, len(modelDeductions))
	totalAmount := 0.0
	for i, md := range modelDeductions {
		deductions[i] = convertModelDeductionToService(&md)
		totalAmount += md.Amount
	}

	return deductions, totalAmount, nil
}

// GetStudentDeductionsWithPagination retrieves all deductions for a student with pagination
func (ds *DeductionService) GetStudentDeductionsWithPagination(studentID uuid.UUID, limit int, offset int) ([]*Deduction, float64, error) {
	var modelDeductions []models.Deduction
	var total int64

	if err := ds.db.Model(&models.Deduction{}).Where("student_id = ?", studentID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count deductions: %w", err)
	}

	if err := ds.db.
		Where("student_id = ?", studentID).
		Order("deduction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&modelDeductions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch deductions: %w", err)
	}

	deductions := make([]*Deduction, len(modelDeductions))
	totalAmount := 0.0
	for i, md := range modelDeductions {
		deductions[i] = convertModelDeductionToService(&md)
		totalAmount += md.Amount
	}

	return deductions, totalAmount, nil
}
