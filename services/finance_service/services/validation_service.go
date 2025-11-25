package services

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"finance_service/database"
	"finance_service/models"
)

// ValidationService handles cross-service validation for stipends and deductions
type ValidationService struct {
	studentClient *StudentServiceClient
	bankingClient *BankingServiceClient
	userClient    *UserServiceClient
	db            *gorm.DB
}

// NewValidationService creates a new validation service
func NewValidationService() *ValidationService {
	return &ValidationService{
		studentClient: NewStudentServiceClient(),
		bankingClient: NewBankingServiceClient(),
		userClient:    NewUserServiceClient(),
		db:            database.DB,
	}
}

// ValidationResult contains validation result details
type ValidationResult struct {
	IsValid  bool
	Errors   []string
	Warnings []string
}

// ============================================================================
// AMOUNT VALIDATION
// ============================================================================

// ValidateAmount checks if an amount is valid (non-negative and within reasonable bounds)
func (vs *ValidationService) ValidateAmount(amount float64, fieldName string) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Check for negative amount
	if amount < 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s cannot be negative: %.2f", fieldName, amount))
	}

	// Check for zero amount (depends on context, so we warn)
	if amount == 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("%s is zero, which may not be intended", fieldName))
	}

	// Check for unreasonably large amounts (100 million as threshold)
	if amount > 100_000_000 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("%s is very large (%.2f), please verify", fieldName, amount))
	}

	return result
}

// ValidateAmountRange checks if an amount falls within a specified range
func (vs *ValidationService) ValidateAmountRange(amount float64, minAmount float64, maxAmount float64, fieldName string) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// First validate the amount itself
	amountValidation := vs.ValidateAmount(amount, fieldName)
	if !amountValidation.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, amountValidation.Errors...)
	}
	result.Warnings = append(result.Warnings, amountValidation.Warnings...)

	// Check minimum
	if amount < minAmount {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("%s (%.2f) is below minimum allowed (%.2f)", fieldName, amount, minAmount))
	}

	// Check maximum
	if amount > maxAmount {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("%s (%.2f) exceeds maximum allowed (%.2f)", fieldName, amount, maxAmount))
	}

	return result
}

// ============================================================================
// DEDUCTION VALIDATION
// ============================================================================

// ValidateDeductionAmount validates a deduction amount against a deduction rule
func (vs *ValidationService) ValidateDeductionAmount(
	deductionAmount float64,
	ruleID uuid.UUID,
	fieldName string,
) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Validate amount is non-negative
	if deductionAmount < 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("%s cannot be negative: %.2f", fieldName, deductionAmount))
		return result
	}

	// Fetch the rule to check limits
	var rule models.DeductionRule
	if err := vs.db.First(&rule, "id = ?", ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Deduction rule not found: %s", ruleID))
		} else {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Error validating rule: %v", err))
		}
		return result
	}

	// Check if rule is active
	if !rule.IsActive {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Deduction rule '%s' is not active", rule.RuleName))
		return result
	}

	// Validate against minimum deduction
	if deductionAmount < rule.MinDeductionAmount {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Deduction amount (%.2f) is below minimum for rule '%s' (%.2f)",
				deductionAmount, rule.RuleName, rule.MinDeductionAmount))
	}

	// Validate against maximum deduction cap
	if deductionAmount > rule.MaxDeductionAmount {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Deduction amount (%.2f) exceeds maximum for rule '%s' (%.2f)",
				deductionAmount, rule.RuleName, rule.MaxDeductionAmount))
	}

	return result
}

// ValidateTotalDeductionAgainstStipend validates that total deductions don't exceed stipend
func (vs *ValidationService) ValidateTotalDeductionAgainstStipend(
	totalDeduction float64,
	stipendAmount float64,
	allowExceed bool,
) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Validate both amounts
	stipendValidation := vs.ValidateAmount(stipendAmount, "Stipend amount")
	deductionValidation := vs.ValidateAmount(totalDeduction, "Total deduction")

	if !stipendValidation.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, stipendValidation.Errors...)
	}
	if !deductionValidation.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, deductionValidation.Errors...)
	}

	if !result.IsValid {
		return result
	}

	// Check if deduction exceeds stipend
	if totalDeduction > stipendAmount {
		if allowExceed {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Total deductions (%.2f) exceed stipend (%.2f), net amount will be negative",
					totalDeduction, stipendAmount))
		} else {
			result.IsValid = false
			result.Errors = append(result.Errors,
				fmt.Sprintf("Total deductions (%.2f) cannot exceed stipend (%.2f)",
					totalDeduction, stipendAmount))
		}
	}

	// Warn if deductions consume more than 80% of stipend
	if stipendAmount > 0 {
		deductionPercentage := (totalDeduction / stipendAmount) * 100
		if deductionPercentage > 80 {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Deductions consume %.1f%% of stipend amount", deductionPercentage))
		}
	}

	return result
}

// ============================================================================
// DEDUCTION RULE VALIDATION
// ============================================================================

// ValidateDeductionRuleInput validates input for creating/updating a deduction rule
func (vs *ValidationService) ValidateDeductionRuleInput(
	ruleName string,
	deductionType string,
	baseAmount float64,
	minAmount float64,
	maxAmount float64,
) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Validate rule name
	if ruleName == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Rule name is required")
	}

	if len(ruleName) > 100 {
		result.IsValid = false
		result.Errors = append(result.Errors, "Rule name exceeds maximum length of 100 characters")
	}

	// Validate deduction type
	if deductionType == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Deduction type is required")
	}

	// Validate base amount
	baseValidation := vs.ValidateAmount(baseAmount, "Base amount")
	if !baseValidation.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, baseValidation.Errors...)
	}

	// Validate min amount
	minValidation := vs.ValidateAmount(minAmount, "Minimum amount")
	if !minValidation.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, minValidation.Errors...)
	}

	// Validate max amount
	maxValidation := vs.ValidateAmount(maxAmount, "Maximum amount")
	if !maxValidation.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, maxValidation.Errors...)
	}

	// Validate amount relationships
	if minAmount > maxAmount {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Minimum amount (%.2f) cannot be greater than maximum amount (%.2f)", minAmount, maxAmount))
	}

	if baseAmount > maxAmount {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Base amount (%.2f) cannot exceed maximum amount (%.2f)", baseAmount, maxAmount))
	}

	if baseAmount < minAmount {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("Base amount (%.2f) is less than minimum (%.2f)", baseAmount, minAmount))
	}

	return result
}

// ============================================================================
// STIPEND VALIDATION
// ============================================================================

// ValidateStipendInput validates input for creating a stipend
func (vs *ValidationService) ValidateStipendInput(
	studentID uuid.UUID,
	stipendType string,
	amount float64,
	journalNumber string,
) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Validate student ID
	if studentID == uuid.Nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "Student ID is required and cannot be nil")
	}

	// Validate stipend type
	if stipendType == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Stipend type is required")
	}

	validTypes := map[string]bool{
		"full-scholarship": true,
		"self-funded":      true,
		"partial":          true,
	}

	if _, valid := validTypes[stipendType]; !valid {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Invalid stipend type '%s'. Must be one of: full-scholarship, self-funded, partial", stipendType))
	}

	// Validate amount
	amountValidation := vs.ValidateAmount(amount, "Stipend amount")
	if !amountValidation.IsValid {
		result.IsValid = false
		result.Errors = append(result.Errors, amountValidation.Errors...)
	}
	result.Warnings = append(result.Warnings, amountValidation.Warnings...)

	// Check upper bound for stipend
	if amount > 10_000_000 {
		result.IsValid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Stipend amount (%.2f) exceeds maximum allowed (10,000,000)", amount))
	}

	// Validate journal number if provided
	if journalNumber != "" {
		if len(journalNumber) > 255 {
			result.IsValid = false
			result.Errors = append(result.Errors, "Journal number exceeds maximum length of 255 characters")
		}

		// Check if journal number already exists (only if database is connected)
		if vs.db != nil {
			var count int64
			if err := vs.db.Model(&models.Stipend{}).Where("journal_number = ?", journalNumber).Count(&count).Error; err == nil && count > 0 {
				result.IsValid = false
				result.Errors = append(result.Errors, fmt.Sprintf("Journal number '%s' already exists", journalNumber))
			}
		}
	}

	return result
}

// ValidateStipendCreation validates all prerequisites for creating a stipend
func (vs *ValidationService) ValidateStipendCreation(userID, studentID string, amount float64) error {
	log.Printf("Validating stipend creation for student %s by user %s with amount %.2f", studentID, userID, amount)

	// 1. Validate User
	log.Println("Step 1: Validating user...")
	if err := vs.userClient.ValidateUserExists(userID); err != nil {
		return fmt.Errorf("user validation failed: %w", err)
	}

	// 2. Validate User Has Permission
	log.Println("Step 2: Validating user permissions...")
	role, err := vs.userClient.ValidateUserPermission(userID)
	if err != nil {
		return fmt.Errorf("user permission validation failed: %w", err)
	}
	log.Printf("User has role: %s", role)

	// 3. Validate Student
	log.Println("Step 3: Validating student...")
	if err := vs.studentClient.ValidateStudent(studentID); err != nil {
		return fmt.Errorf("student validation failed: %w", err)
	}

	// 4. Validate Student Bank Details
	log.Println("Step 4: Validating student bank details...")
	if err := vs.bankingClient.ValidateStudentBankDetails(studentID); err != nil {
		return fmt.Errorf("bank details validation failed: %w", err)
	}

	// 5. Validate Amount
	log.Println("Step 5: Validating amount...")
	amountValidation := vs.ValidateAmount(amount, "Stipend amount")
	if !amountValidation.IsValid {
		return fmt.Errorf("%s", vs.FormatValidationError(amountValidation))
	}

	// Check upper bound
	if amount > 1000000 {
		return fmt.Errorf("stipend amount exceeds maximum allowed")
	}

	log.Println("All validations passed successfully")
	return nil
}

// ValidateDeductionRule validates a deduction rule for applicability
func (vs *ValidationService) ValidateDeductionRule(studentID string, stipendType string, ruleAmount, maxDeduction float64) error {
	// Validate amounts
	if ruleAmount < 0 {
		return fmt.Errorf("deduction amount cannot be negative")
	}

	if maxDeduction < 0 {
		return fmt.Errorf("max deduction amount cannot be negative")
	}

	if ruleAmount > maxDeduction {
		return fmt.Errorf("deduction amount %.2f exceeds max allowed %.2f", ruleAmount, maxDeduction)
	}

	// Validate deduction doesn't exceed stipend
	// This would be checked in database constraints as well
	return nil
}

// GetStudentInfo fetches complete student information for stipend processing
func (vs *ValidationService) GetStudentInfo(studentID string) (*StudentInfo, error) {
	student, err := vs.studentClient.GetStudent(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	bankDetails, err := vs.bankingClient.GetStudentBankDetails(studentID)
	if err != nil {
		// Bank details might not exist yet, but we should know
		log.Printf("Warning: Bank details not available for student %s: %v", studentID, err)
	}

	return &StudentInfo{
		Student:     student,
		BankDetails: bankDetails,
	}, nil
}

// StudentInfo contains consolidated student information from multiple services
type StudentInfo struct {
	Student     *Student
	BankDetails *StudentBankDetails
}

// GetUserInfo fetches complete user information for audit trail
func (vs *ValidationService) GetUserInfo(userID string) (*UserInfo, error) {
	user, err := vs.userClient.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	role, err := vs.userClient.GetRole(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &UserInfo{
		User: user,
		Role: role,
	}, nil
}

// UserInfo contains consolidated user information
type UserInfo struct {
	User *User
	Role *Role
}

// ValidateDeductionApplicability checks if a deduction applies to a student based on stipend type
func (vs *ValidationService) ValidateDeductionApplicability(stipendType string, isApplicableToFullScholar, isApplicableToSelfFunded bool) error {
	switch stipendType {
	case "full-scholarship":
		if !isApplicableToFullScholar {
			return fmt.Errorf("this deduction does not apply to full-scholarship students")
		}
	case "self-funded":
		if !isApplicableToSelfFunded {
			return fmt.Errorf("this deduction does not apply to self-funded students")
		}
	default:
		return fmt.Errorf("unknown stipend type: %s", stipendType)
	}

	return nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// LogValidationResult logs the validation result
func (vs *ValidationService) LogValidationResult(operation string, result *ValidationResult) {
	if result.IsValid {
		log.Printf("[VALIDATION SUCCESS] %s - All validations passed", operation)
	} else {
		log.Printf("[VALIDATION FAILED] %s - Errors: %v", operation, result.Errors)
	}

	if len(result.Warnings) > 0 {
		log.Printf("[VALIDATION WARNING] %s - Warnings: %v", operation, result.Warnings)
	}
}

// FormatValidationError formats validation errors into a user-friendly message
func (vs *ValidationService) FormatValidationError(result *ValidationResult) string {
	if result.IsValid {
		return ""
	}

	errorMsg := "Validation failed:\n"
	for i, err := range result.Errors {
		errorMsg += fmt.Sprintf("  %d. %s\n", i+1, err)
	}

	return errorMsg
}

// FormatValidationWarnings formats validation warnings into a user-friendly message
func (vs *ValidationService) FormatValidationWarnings(result *ValidationResult) string {
	if len(result.Warnings) == 0 {
		return ""
	}

	warningMsg := "Validation warnings:\n"
	for i, warn := range result.Warnings {
		warningMsg += fmt.Sprintf("  %d. %s\n", i+1, warn)
	}

	return warningMsg
}

