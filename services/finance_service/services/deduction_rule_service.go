package services

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"finance_service/database"
	"finance_service/models"
)

// DeductionRuleService handles business logic for deduction rules management
type DeductionRuleService struct {
	db     *gorm.DB
	logger *ErrorLogger
}

// NewDeductionRuleService creates a new deduction rule service
func NewDeductionRuleService() *DeductionRuleService {
	return &DeductionRuleService{
		db:     database.DB,
		logger: NewErrorLogger(1000),
	}
}

// CreateDeductionRuleInput represents input for creating a deduction rule
type CreateDeductionRuleInput struct {
	RuleName                  string
	DeductionType             string
	Description               string
	BaseAmount                float64
	MaxDeductionAmount        float64
	MinDeductionAmount        float64
	IsApplicableToFullScholar bool
	IsApplicableToSelfFunded  bool
	AppliesMonthly            bool
	AppliesAnnually           bool
	IsOptional                bool
	Priority                  int
	CreatedBy                 *uuid.UUID
}

// UpdateDeductionRuleInput represents input for updating a deduction rule
type UpdateDeductionRuleInput struct {
	RuleName                  *string
	DeductionType             *string
	Description               *string
	BaseAmount                *float64
	MaxDeductionAmount        *float64
	MinDeductionAmount        *float64
	IsApplicableToFullScholar *bool
	IsApplicableToSelfFunded  *bool
	AppliesMonthly            *bool
	AppliesAnnually           *bool
	IsOptional                *bool
	Priority                  *int
	IsActive                  *bool
	ModifiedBy                *uuid.UUID
}

// CreateRule creates a new deduction rule with comprehensive validation
func (drs *DeductionRuleService) CreateRule(input *CreateDeductionRuleInput) (*models.DeductionRule, error) {
	log.Printf("Creating deduction rule: %s", input.RuleName)

	// Validation
	if err := drs.validateRuleInput(input); err != nil {
		drs.logger.LogError(CategoryDeductionValidation, "Deduction rule creation validation failed", 
			map[string]string{"error": err.Error()})
		return nil, err
	}

	// Check for duplicate rule name
	var existing models.DeductionRule
	if err := drs.db.Where("rule_name = ?", input.RuleName).First(&existing).Error; err == nil {
		errMsg := fmt.Sprintf("rule name '%s' already exists", input.RuleName)
		drs.logger.LogError(CategoryDeductionValidation, errMsg, nil)
		return nil, fmt.Errorf(errMsg)
	} else if err != gorm.ErrRecordNotFound {
		drs.logger.LogError(CategoryDatabaseError, "Failed to check duplicate rule name", 
			map[string]string{"error": err.Error()})
		return nil, fmt.Errorf("failed to check duplicate rule name: %w", err)
	}

	// Create rule
	rule := &models.DeductionRule{
		ID:                        uuid.New(),
		RuleName:                  input.RuleName,
		DeductionType:             input.DeductionType,
		Description:               input.Description,
		BaseAmount:                input.BaseAmount,
		MaxDeductionAmount:        input.MaxDeductionAmount,
		MinDeductionAmount:        input.MinDeductionAmount,
		IsApplicableToFullScholar: input.IsApplicableToFullScholar,
		IsApplicableToSelfFunded:  input.IsApplicableToSelfFunded,
		AppliesMonthly:            input.AppliesMonthly,
		AppliesAnnually:           input.AppliesAnnually,
		IsOptional:                input.IsOptional,
		Priority:                  input.Priority,
		IsActive:                  true,
		CreatedBy:                 input.CreatedBy,
	}

	if err := drs.db.Create(rule).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to create deduction rule", 
			map[string]string{"rule_name": input.RuleName, "error": err.Error()})
		return nil, fmt.Errorf("failed to create deduction rule: %w", err)
	}

	drs.logger.LogInfo(CategoryDeductionValidation, 
		fmt.Sprintf("Deduction rule created: %s", rule.RuleName), 
		map[string]string{"rule_id": rule.ID.String()})

	log.Printf("Deduction rule created with ID: %s", rule.ID)
	return rule, nil
}

// GetRuleByID retrieves a deduction rule by ID
func (drs *DeductionRuleService) GetRuleByID(ruleID uuid.UUID) (*models.DeductionRule, error) {
	var rule models.DeductionRule

	if err := drs.db.First(&rule, "id = ?", ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			drs.logger.LogWarning(CategoryDeductionValidation, 
				fmt.Sprintf("Deduction rule not found: %s", ruleID), nil)
			return nil, fmt.Errorf("deduction rule not found")
		}
		drs.logger.LogError(CategoryDatabaseError, "Failed to fetch deduction rule", 
			map[string]string{"rule_id": ruleID.String(), "error": err.Error()})
		return nil, fmt.Errorf("failed to fetch deduction rule: %w", err)
	}

	return &rule, nil
}

// ListActiveRules retrieves all active deduction rules with pagination
func (drs *DeductionRuleService) ListActiveRules(limit int, offset int) ([]models.DeductionRule, int64, error) {
	var rules []models.DeductionRule
	var total int64

	if err := drs.db.Model(&models.DeductionRule{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to count active deduction rules", 
			map[string]string{"error": err.Error()})
		return nil, 0, fmt.Errorf("failed to count rules: %w", err)
	}

	if err := drs.db.
		Where("is_active = ?", true).
		Order("priority DESC, rule_name ASC").
		Limit(limit).
		Offset(offset).
		Find(&rules).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to fetch active deduction rules", 
			map[string]string{"error": err.Error()})
		return nil, 0, fmt.Errorf("failed to fetch rules: %w", err)
	}

	return rules, total, nil
}

// ListAllRules retrieves all deduction rules (active and inactive) with pagination
func (drs *DeductionRuleService) ListAllRules(limit int, offset int) ([]models.DeductionRule, int64, error) {
	var rules []models.DeductionRule
	var total int64

	if err := drs.db.Model(&models.DeductionRule{}).Count(&total).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to count deduction rules", 
			map[string]string{"error": err.Error()})
		return nil, 0, fmt.Errorf("failed to count rules: %w", err)
	}

	if err := drs.db.
		Order("priority DESC, rule_name ASC").
		Limit(limit).
		Offset(offset).
		Find(&rules).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to fetch deduction rules", 
			map[string]string{"error": err.Error()})
		return nil, 0, fmt.Errorf("failed to fetch rules: %w", err)
	}

	return rules, total, nil
}

// ListRulesByType retrieves deduction rules filtered by type
func (drs *DeductionRuleService) ListRulesByType(deductionType string, limit int, offset int) ([]models.DeductionRule, int64, error) {
	var rules []models.DeductionRule
	var total int64

	if deductionType == "" {
		return nil, 0, fmt.Errorf("deduction type is required")
	}

	if err := drs.db.Model(&models.DeductionRule{}).
		Where("deduction_type = ? AND is_active = ?", deductionType, true).
		Count(&total).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to count rules by type", 
			map[string]string{"deduction_type": deductionType, "error": err.Error()})
		return nil, 0, fmt.Errorf("failed to count rules: %w", err)
	}

	if err := drs.db.
		Where("deduction_type = ? AND is_active = ?", deductionType, true).
		Order("priority DESC, rule_name ASC").
		Limit(limit).
		Offset(offset).
		Find(&rules).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to fetch rules by type", 
			map[string]string{"deduction_type": deductionType, "error": err.Error()})
		return nil, 0, fmt.Errorf("failed to fetch rules: %w", err)
	}

	return rules, total, nil
}

// UpdateRule updates an existing deduction rule with validation
func (drs *DeductionRuleService) UpdateRule(ruleID uuid.UUID, input *UpdateDeductionRuleInput) (*models.DeductionRule, error) {
	log.Printf("Updating deduction rule: %s", ruleID)

	// Fetch existing rule
	rule, err := drs.GetRuleByID(ruleID)
	if err != nil {
		return nil, err
	}

	// Validate update input
	if err := drs.validateRuleUpdate(input); err != nil {
		drs.logger.LogError(CategoryDeductionValidation, "Deduction rule update validation failed", 
			map[string]string{"rule_id": ruleID.String(), "error": err.Error()})
		return nil, err
	}

	// Check for duplicate rule name if updating
	if input.RuleName != nil && *input.RuleName != rule.RuleName {
		var existing models.DeductionRule
		if err := drs.db.Where("rule_name = ? AND id != ?", *input.RuleName, ruleID).First(&existing).Error; err == nil {
			errMsg := fmt.Sprintf("rule name '%s' already exists", *input.RuleName)
			drs.logger.LogError(CategoryDeductionValidation, errMsg, nil)
			return nil, fmt.Errorf(errMsg)
		} else if err != gorm.ErrRecordNotFound {
			drs.logger.LogError(CategoryDatabaseError, "Failed to check duplicate rule name", 
				map[string]string{"error": err.Error()})
			return nil, fmt.Errorf("failed to check duplicate rule name: %w", err)
		}
	}

	// Apply updates
	if input.RuleName != nil {
		rule.RuleName = *input.RuleName
	}
	if input.DeductionType != nil {
		rule.DeductionType = *input.DeductionType
	}
	if input.Description != nil {
		rule.Description = *input.Description
	}
	if input.BaseAmount != nil {
		rule.BaseAmount = *input.BaseAmount
	}
	if input.MaxDeductionAmount != nil {
		rule.MaxDeductionAmount = *input.MaxDeductionAmount
	}
	if input.MinDeductionAmount != nil {
		rule.MinDeductionAmount = *input.MinDeductionAmount
	}
	if input.IsApplicableToFullScholar != nil {
		rule.IsApplicableToFullScholar = *input.IsApplicableToFullScholar
	}
	if input.IsApplicableToSelfFunded != nil {
		rule.IsApplicableToSelfFunded = *input.IsApplicableToSelfFunded
	}
	if input.AppliesMonthly != nil {
		rule.AppliesMonthly = *input.AppliesMonthly
	}
	if input.AppliesAnnually != nil {
		rule.AppliesAnnually = *input.AppliesAnnually
	}
	if input.IsOptional != nil {
		rule.IsOptional = *input.IsOptional
	}
	if input.Priority != nil {
		rule.Priority = *input.Priority
	}
	if input.IsActive != nil {
		rule.IsActive = *input.IsActive
	}
	if input.ModifiedBy != nil {
		rule.ModifiedBy = input.ModifiedBy
	}

	// Save updates
	if err := drs.db.Save(rule).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to update deduction rule", 
			map[string]string{"rule_id": ruleID.String(), "error": err.Error()})
		return nil, fmt.Errorf("failed to update deduction rule: %w", err)
	}

	drs.logger.LogInfo(CategoryDeductionValidation, 
		fmt.Sprintf("Deduction rule updated: %s", rule.RuleName), 
		map[string]string{"rule_id": ruleID.String()})

	log.Printf("Deduction rule %s updated successfully", ruleID)
	return rule, nil
}

// DeleteRule soft-deletes a deduction rule by marking it as inactive
func (drs *DeductionRuleService) DeleteRule(ruleID uuid.UUID) error {
	log.Printf("Deleting deduction rule: %s", ruleID)

	rule, err := drs.GetRuleByID(ruleID)
	if err != nil {
		return err
	}

	// Soft delete by marking as inactive
	if err := drs.db.Model(rule).Update("is_active", false).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to delete deduction rule", 
			map[string]string{"rule_id": ruleID.String(), "error": err.Error()})
		return fmt.Errorf("failed to delete deduction rule: %w", err)
	}

	drs.logger.LogInfo(CategoryDeductionValidation, 
		fmt.Sprintf("Deduction rule deleted: %s", rule.RuleName), 
		map[string]string{"rule_id": ruleID.String()})

	log.Printf("Deduction rule %s deleted successfully", ruleID)
	return nil
}

// GetApplicableRules retrieves rules applicable to a specific student type
func (drs *DeductionRuleService) GetApplicableRules(isFullScholar bool) ([]models.DeductionRule, error) {
	var rules []models.DeductionRule

	query := drs.db.Where("is_active = ?", true)
	if isFullScholar {
		query = query.Where("is_applicable_to_full_scholar = ?", true)
	} else {
		query = query.Where("is_applicable_to_self_funded = ?", true)
	}

	if err := query.Order("priority DESC, rule_name ASC").Find(&rules).Error; err != nil {
		drs.logger.LogError(CategoryDatabaseError, "Failed to fetch applicable rules", 
			map[string]string{"error": err.Error()})
		return nil, fmt.Errorf("failed to fetch applicable rules: %w", err)
	}

	return rules, nil
}

// ============================================================================
// VALIDATION HELPERS
// ============================================================================

// validateRuleInput validates the input for creating a deduction rule
func (drs *DeductionRuleService) validateRuleInput(input *CreateDeductionRuleInput) error {
	if input.RuleName == "" {
		return fmt.Errorf("rule name is required")
	}

	if len(input.RuleName) > 100 {
		return fmt.Errorf("rule name cannot exceed 100 characters")
	}

	if input.DeductionType == "" {
		return fmt.Errorf("deduction type is required")
	}

	if input.BaseAmount < 0 {
		return fmt.Errorf("base amount cannot be negative")
	}

	if input.MaxDeductionAmount < 0 {
		return fmt.Errorf("max deduction amount cannot be negative")
	}

	if input.MinDeductionAmount < 0 {
		return fmt.Errorf("min deduction amount cannot be negative")
	}

	if input.MaxDeductionAmount < input.MinDeductionAmount {
		return fmt.Errorf("max deduction amount must be >= min deduction amount")
	}

	if !input.AppliesMonthly && !input.AppliesAnnually {
		return fmt.Errorf("rule must apply either monthly or annually")
	}

	return nil
}

// validateRuleUpdate validates the input for updating a deduction rule
func (drs *DeductionRuleService) validateRuleUpdate(input *UpdateDeductionRuleInput) error {
	if input.RuleName != nil && *input.RuleName == "" {
		return fmt.Errorf("rule name cannot be empty")
	}

	if input.RuleName != nil && len(*input.RuleName) > 100 {
		return fmt.Errorf("rule name cannot exceed 100 characters")
	}

	if input.DeductionType != nil && *input.DeductionType == "" {
		return fmt.Errorf("deduction type cannot be empty")
	}

	if input.BaseAmount != nil && *input.BaseAmount < 0 {
		return fmt.Errorf("base amount cannot be negative")
	}

	if input.MaxDeductionAmount != nil && *input.MaxDeductionAmount < 0 {
		return fmt.Errorf("max deduction amount cannot be negative")
	}

	if input.MinDeductionAmount != nil && *input.MinDeductionAmount < 0 {
		return fmt.Errorf("min deduction amount cannot be negative")
	}

	return nil
}

// GetErrorLogger returns the error logger for logging purposes
func (drs *DeductionRuleService) GetErrorLogger() *ErrorLogger {
	return drs.logger
}
