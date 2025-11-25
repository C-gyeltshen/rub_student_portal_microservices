package services

// Validation constants and limits for financial operations
const (
	// Amount limits
	MinAllowedAmount           = 0.0
	MaxStipendAmount           = 10_000_000.0  // 10 million max stipend
	MaxDeductionAmount         = 100_000_000.0 // 100 million max deduction (very large threshold warning)
	WarningThresholdAmount     = 100_000_000.0 // Warn if amount exceeds this
	DeductionPercentageWarning = 80.0          // Warn if deductions exceed 80% of stipend

	// Rule name limits
	MaxRuleNameLength    = 100
	MaxJournalNumberLen  = 255
	MaxDeductionTypeLen  = 100
	MaxDescriptionLen    = 5000

	// Deduction thresholds
	MaxTotalDeductionPercentage = 100.0 // Allow 100% deduction at database level, warn at service level

	// Valid stipend types
	StipendTypeFullScholarship = "full-scholarship"
	StipendTypeSelfFunded      = "self-funded"
	StipendTypePartial         = "partial"

	// Deduction application frequencies
	FrequencyMonthly = "Monthly"
	FrequencyAnnual  = "Annual"

	// Applicability flags
	ApplicableToAll            = "All"
	ApplicableToFullScholarship = "FullScholarship"
	ApplicableToSelfFunded     = "SelfFunded"

	// Processing statuses
	StatusPending  = "Pending"
	StatusApproved = "Approved"
	StatusProcessed = "Processed"
	StatusRejected = "Rejected"

	// Payment statuses
	PaymentStatusPending  = "Pending"
	PaymentStatusProcessed = "Processed"
	PaymentStatusFailed    = "Failed"
)

// ValidationLimits defines all validation limits for financial operations
type ValidationLimits struct {
	// Stipend limits
	MinStipendAmount float64
	MaxStipendAmount float64

	// Deduction limits
	MinDeductionAmount     float64
	MaxDeductionAmount     float64
	MaxDeductionPercentage float64

	// Rule limits
	MaxRuleNameLength int
	MaxDescriptionLen int

	// Warning thresholds
	DeductionPercentageWarning float64
	LargeAmountThreshold       float64
}

// DefaultValidationLimits returns the default validation limits
func DefaultValidationLimits() ValidationLimits {
	return ValidationLimits{
		MinStipendAmount:           0.0,
		MaxStipendAmount:           10_000_000.0,
		MinDeductionAmount:         0.0,
		MaxDeductionAmount:         100_000_000.0,
		MaxDeductionPercentage:     100.0,
		MaxRuleNameLength:          100,
		MaxDescriptionLen:          5000,
		DeductionPercentageWarning: 80.0,
		LargeAmountThreshold:       100_000_000.0,
	}
}

// ValidStipendTypes returns the list of valid stipend types
func ValidStipendTypes() []string {
	return []string{
		StipendTypeFullScholarship,
		StipendTypeSelfFunded,
		StipendTypePartial,
	}
}

// ValidFrequencies returns the list of valid deduction frequencies
func ValidFrequencies() []string {
	return []string{
		FrequencyMonthly,
		FrequencyAnnual,
	}
}

// ValidApplicabilities returns the list of valid applicability flags
func ValidApplicabilities() []string {
	return []string{
		ApplicableToAll,
		ApplicableToFullScholarship,
		ApplicableToSelfFunded,
	}
}

// ValidProcessingStatuses returns the list of valid processing statuses
func ValidProcessingStatuses() []string {
	return []string{
		StatusPending,
		StatusApproved,
		StatusProcessed,
		StatusRejected,
	}
}

// ValidPaymentStatuses returns the list of valid payment statuses
func ValidPaymentStatuses() []string {
	return []string{
		PaymentStatusPending,
		PaymentStatusProcessed,
		PaymentStatusFailed,
	}
}
