package services

import (
	"testing"

	"github.com/google/uuid"
)

// TestValidateAmount tests the ValidateAmount function
func TestValidateAmount(t *testing.T) {
	vs := NewValidationService()

	tests := []struct {
		name        string
		amount      float64
		fieldName   string
		shouldValid bool
		hasWarning  bool
	}{
		{
			name:        "valid positive amount",
			amount:      5000.50,
			fieldName:   "Amount",
			shouldValid: true,
			hasWarning:  false,
		},
		{
			name:        "negative amount",
			amount:      -100.0,
			fieldName:   "Amount",
			shouldValid: false,
			hasWarning:  false,
		},
		{
			name:        "zero amount",
			amount:      0.0,
			fieldName:   "Amount",
			shouldValid: true,
			hasWarning:  true, // warning for zero
		},
		{
			name:        "very large amount",
			amount:      200_000_000.0,
			fieldName:   "Amount",
			shouldValid: true,
			hasWarning:  true, // warning for large amount
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.ValidateAmount(tt.amount, tt.fieldName)

			if result.IsValid != tt.shouldValid {
				t.Errorf("expected IsValid=%v, got %v", tt.shouldValid, result.IsValid)
			}

			hasWarning := len(result.Warnings) > 0
			if hasWarning != tt.hasWarning {
				t.Errorf("expected hasWarning=%v, got %v. Warnings: %v", tt.hasWarning, hasWarning, result.Warnings)
			}
		})
	}
}

// TestValidateAmountRange tests the ValidateAmountRange function
func TestValidateAmountRange(t *testing.T) {
	vs := NewValidationService()

	tests := []struct {
		name        string
		amount      float64
		minAmount   float64
		maxAmount   float64
		shouldValid bool
	}{
		{
			name:        "amount within range",
			amount:      5000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: true,
		},
		{
			name:        "amount below minimum",
			amount:      500.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: false,
		},
		{
			name:        "amount above maximum",
			amount:      15000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: false,
		},
		{
			name:        "negative amount",
			amount:      -5000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.ValidateAmountRange(tt.amount, tt.minAmount, tt.maxAmount, "Test Amount")

			if result.IsValid != tt.shouldValid {
				t.Errorf("expected IsValid=%v, got %v. Errors: %v", tt.shouldValid, result.IsValid, result.Errors)
			}
		})
	}
}

// TestValidateDeductionRuleInput tests deduction rule validation
func TestValidateDeductionRuleInput(t *testing.T) {
	vs := NewValidationService()

	tests := []struct {
		name        string
		ruleName    string
		deductType  string
		baseAmount  float64
		minAmount   float64
		maxAmount   float64
		shouldValid bool
	}{
		{
			name:        "valid rule",
			ruleName:    "Hostel Fee",
			deductType:  "hostel",
			baseAmount:  5000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: true,
		},
		{
			name:        "empty rule name",
			ruleName:    "",
			deductType:  "hostel",
			baseAmount:  5000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: false,
		},
		{
			name:        "empty deduction type",
			ruleName:    "Hostel Fee",
			deductType:  "",
			baseAmount:  5000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: false,
		},
		{
			name:        "negative base amount",
			ruleName:    "Hostel Fee",
			deductType:  "hostel",
			baseAmount:  -5000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: false,
		},
		{
			name:        "min > max",
			ruleName:    "Hostel Fee",
			deductType:  "hostel",
			baseAmount:  5000.0,
			minAmount:   10000.0,
			maxAmount:   5000.0,
			shouldValid: false,
		},
		{
			name:        "base > max",
			ruleName:    "Hostel Fee",
			deductType:  "hostel",
			baseAmount:  15000.0,
			minAmount:   1000.0,
			maxAmount:   10000.0,
			shouldValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.ValidateDeductionRuleInput(
				tt.ruleName,
				tt.deductType,
				tt.baseAmount,
				tt.minAmount,
				tt.maxAmount,
			)

			if result.IsValid != tt.shouldValid {
				t.Errorf("expected IsValid=%v, got %v. Errors: %v", tt.shouldValid, result.IsValid, result.Errors)
			}
		})
	}
}

// TestValidateTotalDeductionAgainstStipend tests total deduction validation
func TestValidateTotalDeductionAgainstStipend(t *testing.T) {
	vs := NewValidationService()

	tests := []struct {
		name              string
		totalDeduction    float64
		stipendAmount     float64
		allowExceed       bool
		shouldValid       bool
		shouldHaveWarning bool
	}{
		{
			name:              "deductions within limit",
			totalDeduction:    4000.0,
			stipendAmount:     5000.0,
			allowExceed:       false,
			shouldValid:       true,
			shouldHaveWarning: false,
		},
		{
			name:              "deductions equal stipend",
			totalDeduction:    5000.0,
			stipendAmount:     5000.0,
			allowExceed:       false,
			shouldValid:       true,
			shouldHaveWarning: true, // 100% threshold
		},
		{
			name:              "deductions exceed stipend (not allowed)",
			totalDeduction:    6000.0,
			stipendAmount:     5000.0,
			allowExceed:       false,
			shouldValid:       false,
			shouldHaveWarning: true, // warning for exceeding 80% (which happens even when not allowed)
		},
		{
			name:              "deductions exceed stipend (allowed)",
			totalDeduction:    6000.0,
			stipendAmount:     5000.0,
			allowExceed:       true,
			shouldValid:       true,
			shouldHaveWarning: true, // exceeds 80%
		},
		{
			name:              "high deduction percentage",
			totalDeduction:    4200.0,
			stipendAmount:     5000.0,
			allowExceed:       false,
			shouldValid:       true,
			shouldHaveWarning: true, // 84% > 80%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.ValidateTotalDeductionAgainstStipend(
				tt.totalDeduction,
				tt.stipendAmount,
				tt.allowExceed,
			)

			if result.IsValid != tt.shouldValid {
				t.Errorf("expected IsValid=%v, got %v. Errors: %v", tt.shouldValid, result.IsValid, result.Errors)
			}

			hasWarning := len(result.Warnings) > 0
			if hasWarning != tt.shouldHaveWarning {
				t.Errorf("expected hasWarning=%v, got %v. Warnings: %v", tt.shouldHaveWarning, hasWarning, result.Warnings)
			}
		})
	}
}

// TestValidateStipendInput tests stipend input validation
func TestValidateStipendInput(t *testing.T) {
	vs := NewValidationService()

	validStudentID := uuid.New()
	nilUUID := uuid.Nil

	tests := []struct {
		name        string
		studentID   uuid.UUID
		stipendType string
		amount      float64
		journal     string
		shouldValid bool
	}{
		{
			name:        "valid stipend input",
			studentID:   validStudentID,
			stipendType: "full-scholarship",
			amount:      100000.0,
			journal:     "JN-STI-2024-001",
			shouldValid: true,
		},
		{
			name:        "nil student ID",
			studentID:   nilUUID,
			stipendType: "full-scholarship",
			amount:      100000.0,
			journal:     "JN-STI-2024-001",
			shouldValid: false,
		},
		{
			name:        "invalid stipend type",
			studentID:   validStudentID,
			stipendType: "invalid-type",
			amount:      100000.0,
			journal:     "JN-STI-2024-001",
			shouldValid: false,
		},
		{
			name:        "negative amount",
			studentID:   validStudentID,
			stipendType: "full-scholarship",
			amount:      -100000.0,
			journal:     "JN-STI-2024-001",
			shouldValid: false,
		},
		{
			name:        "amount exceeds maximum",
			studentID:   validStudentID,
			stipendType: "full-scholarship",
			amount:      15_000_000.0,
			journal:     "JN-STI-2024-001",
			shouldValid: false,
		},
		{
			name:        "valid self-funded",
			studentID:   validStudentID,
			stipendType: "self-funded",
			amount:      50000.0,
			journal:     "JN-STI-2024-002",
			shouldValid: true,
		},
		{
			name:        "valid partial",
			studentID:   validStudentID,
			stipendType: "partial",
			amount:      75000.0,
			journal:     "JN-STI-2024-003",
			shouldValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.ValidateStipendInput(
				tt.studentID,
				tt.stipendType,
				tt.amount,
				tt.journal,
			)

			if result.IsValid != tt.shouldValid {
				t.Errorf("expected IsValid=%v, got %v. Errors: %v", tt.shouldValid, result.IsValid, result.Errors)
			}
		})
	}
}

// TestFormatValidationError tests error formatting
func TestFormatValidationError(t *testing.T) {
	vs := NewValidationService()

	result := &ValidationResult{
		IsValid: false,
		Errors: []string{
			"Error 1",
			"Error 2",
		},
		Warnings: []string{},
	}

	errorMsg := vs.FormatValidationError(result)
	if errorMsg == "" {
		t.Errorf("expected non-empty error message")
	}

	if len(errorMsg) == 0 {
		t.Errorf("error message is empty")
	}

	t.Logf("Formatted error:\n%s", errorMsg)
}

// TestFormatValidationWarnings tests warning formatting
func TestFormatValidationWarnings(t *testing.T) {
	vs := NewValidationService()

	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{"Warning 1", "Warning 2"},
	}

	warningMsg := vs.FormatValidationWarnings(result)
	if warningMsg == "" {
		t.Errorf("expected non-empty warning message")
	}

	t.Logf("Formatted warnings:\n%s", warningMsg)
}
