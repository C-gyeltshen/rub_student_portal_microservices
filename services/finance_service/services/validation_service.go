package services

import (
	"fmt"
	"log"
)

// ValidationService handles cross-service validation for stipends and deductions
type ValidationService struct {
	studentClient *StudentServiceClient
	bankingClient *BankingServiceClient
	userClient    *UserServiceClient
}

// NewValidationService creates a new validation service
func NewValidationService() *ValidationService {
	return &ValidationService{
		studentClient: NewStudentServiceClient(),
		bankingClient: NewBankingServiceClient(),
		userClient:    NewUserServiceClient(),
	}
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
	if amount <= 0 {
		return fmt.Errorf("stipend amount must be positive")
	}

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
