package services

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"finance_service/database"
	"finance_service/models"
)

// TestStipendCalculationFullScholarship tests stipend calculation for full-scholarship students
func TestStipendCalculationFullScholarship(t *testing.T) {
	// Setup
	setupTestDBNew(t)
	service := NewStipendService()

	studentID := getTestStudentID()
	baseAmount := 50000.0

	// Create test deduction rules for full-scholarship students
	rules := []models.DeductionRule{
		{
			ID:                        uuid.New(),
			RuleName:                  "Hostel Fee",
			DeductionType:             "hostel",
			BaseAmount:                3000.0,
			MaxDeductionAmount:        3500.0,
			MinDeductionAmount:        2500.0,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  false,
			IsActive:                  true,
			Priority:                  100,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "Electricity",
			DeductionType:             "electricity",
			BaseAmount:                500.0,
			MaxDeductionAmount:        800.0,
			MinDeductionAmount:        300.0,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  false,
			IsActive:                  true,
			Priority:                  50,
		},
	}

	for _, rule := range rules {
		if err := database.DB.Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create test rule: %v", err)
		}
	}

	// Test calculation
	result, err := service.CalculateStipendWithDeductions(studentID, "full-scholarship", baseAmount)
	if err != nil {
		t.Fatalf("CalculateStipendWithDeductions failed: %v", err)
	}

	// Assertions
	expectedTotalDeductions := 3000.0 + 500.0 // Hostel + Electricity
	if result.BaseStipendAmount != baseAmount {
		t.Errorf("Base amount mismatch: expected %.2f, got %.2f", baseAmount, result.BaseStipendAmount)
	}

	if result.TotalDeductions != expectedTotalDeductions {
		t.Errorf("Total deductions mismatch: expected %.2f, got %.2f", expectedTotalDeductions, result.TotalDeductions)
	}

	expectedNet := baseAmount - expectedTotalDeductions
	if result.NetStipendAmount != expectedNet {
		t.Errorf("Net amount mismatch: expected %.2f, got %.2f", expectedNet, result.NetStipendAmount)
	}

	if len(result.Deductions) != 2 {
		t.Errorf("Expected 2 deductions, got %d", len(result.Deductions))
	}

	// Cleanup
	cleanupTestDB(t, rules)
}

// TestStipendCalculationSelfFunded tests stipend calculation for self-funded students
func TestStipendCalculationSelfFunded(t *testing.T) {
	// Setup
	setupTestDBNew(t)
	service := NewStipendService()

	studentID := getTestStudentID()
	baseAmount := 40000.0

	// Create test deduction rules for self-funded students
	rules := []models.DeductionRule{
		{
			ID:                        uuid.New(),
			RuleName:                  "Hostel Fee",
			DeductionType:             "hostel",
			BaseAmount:                3000.0,
			MaxDeductionAmount:        3500.0,
			MinDeductionAmount:        2500.0,
			IsApplicableToFullScholar: false,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			Priority:                  100,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "Mess Fee",
			DeductionType:             "mess_fees",
			BaseAmount:                2000.0,
			MaxDeductionAmount:        2500.0,
			MinDeductionAmount:        1500.0,
			IsApplicableToFullScholar: false,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			Priority:                  90,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "Electricity",
			DeductionType:             "electricity",
			BaseAmount:                800.0,
			MaxDeductionAmount:        1000.0,
			MinDeductionAmount:        500.0,
			IsApplicableToFullScholar: false,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			Priority:                  50,
		},
	}

	for _, rule := range rules {
		if err := database.DB.Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create test rule: %v", err)
		}
	}

	// Test calculation
	result, err := service.CalculateStipendWithDeductions(studentID, "self-funded", baseAmount)
	if err != nil {
		t.Fatalf("CalculateStipendWithDeductions failed: %v", err)
	}

	// Assertions
	expectedTotalDeductions := 3000.0 + 2000.0 + 800.0 // Hostel + Mess + Electricity
	if result.TotalDeductions != expectedTotalDeductions {
		t.Errorf("Total deductions mismatch: expected %.2f, got %.2f", expectedTotalDeductions, result.TotalDeductions)
	}

	expectedNet := baseAmount - expectedTotalDeductions
	if result.NetStipendAmount != expectedNet {
		t.Errorf("Net amount mismatch: expected %.2f, got %.2f", expectedNet, result.NetStipendAmount)
	}

	if len(result.Deductions) != 3 {
		t.Errorf("Expected 3 deductions, got %d", len(result.Deductions))
	}

	// Cleanup
	cleanupTestDB(t, rules)
}

// TestDeductionPriority tests that deductions are applied in priority order
func TestDeductionPriority(t *testing.T) {
	// Setup
	setupTestDBNew(t)
	service := NewStipendService()

	studentID := getTestStudentID()
	baseAmount := 10000.0

	// Create rules with different priorities
	rules := []models.DeductionRule{
		{
			ID:                        uuid.New(),
			RuleName:                  "Low Priority",
			DeductionType:             "test1",
			BaseAmount:                3000.0,
			MaxDeductionAmount:        3000.0,
			MinDeductionAmount:        0,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			Priority:                  1,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "High Priority",
			DeductionType:             "test2",
			BaseAmount:                4000.0,
			MaxDeductionAmount:        4000.0,
			MinDeductionAmount:        0,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			Priority:                  100,
		},
	}

	for _, rule := range rules {
		if err := database.DB.Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create test rule: %v", err)
		}
	}

	// Test calculation
	result, err := service.CalculateStipendWithDeductions(studentID, "self-funded", baseAmount)
	if err != nil {
		t.Fatalf("CalculateStipendWithDeductions failed: %v", err)
	}

	// High priority should be first
	if result.Deductions[0].RuleName != "High Priority" {
		t.Errorf("High priority deduction should be first, got %s", result.Deductions[0].RuleName)
	}

	// Cleanup
	cleanupTestDB(t, rules)
}

// TestDeductionCapToRemainingStipend tests that deductions don't exceed remaining stipend
func TestDeductionCapToRemainingStipend(t *testing.T) {
	// Setup
	setupTestDBNew(t)
	service := NewStipendService()

	studentID := getTestStudentID()
	baseAmount := 5000.0 // Small amount

	// Create rules that would exceed stipend if applied fully
	rules := []models.DeductionRule{
		{
			ID:                        uuid.New(),
			RuleName:                  "Deduction 1",
			DeductionType:             "test1",
			BaseAmount:                3000.0,
			MaxDeductionAmount:        3000.0,
			MinDeductionAmount:        0,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			Priority:                  100,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "Deduction 2",
			DeductionType:             "test2",
			BaseAmount:                3000.0,
			MaxDeductionAmount:        3000.0,
			MinDeductionAmount:        0,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			Priority:                  50,
		},
	}

	for _, rule := range rules {
		if err := database.DB.Create(&rule).Error; err != nil {
			t.Fatalf("Failed to create test rule: %v", err)
		}
	}

	// Test calculation
	result, err := service.CalculateStipendWithDeductions(studentID, "self-funded", baseAmount)
	if err != nil {
		t.Fatalf("CalculateStipendWithDeductions failed: %v", err)
	}

	// Net should not be negative
	if result.NetStipendAmount < 0 {
		t.Errorf("Net amount should not be negative, got %.2f", result.NetStipendAmount)
	}

	// Total deductions should not exceed base
	if result.TotalDeductions > baseAmount {
		t.Errorf("Total deductions %.2f exceed base amount %.2f", result.TotalDeductions, baseAmount)
	}

	// Cleanup
	cleanupTestDB(t, rules)
}

// TestMonthlyStipendCalculation tests monthly stipend calculation
func TestMonthlyStipendCalculation(t *testing.T) {
	// Setup
	setupTestDBNew(t)
	service := NewStipendService()

	studentID := getTestStudentID()
	annualAmount := 600000.0
	expectedMonthly := annualAmount / 12 // 50,000

	// Create a simple deduction rule
	rule := models.DeductionRule{
		ID:                        uuid.New(),
		RuleName:                  "Test Deduction",
		DeductionType:             "test",
		BaseAmount:                5000.0,
		MaxDeductionAmount:        5000.0,
		MinDeductionAmount:        0,
		IsApplicableToFullScholar: true,
		IsApplicableToSelfFunded:  true,
		IsActive:                  true,
		Priority:                  1,
	}

	if err := database.DB.Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}

	// Test calculation
	result, err := service.CalculateMonthlyStipendForStudent(studentID, "full-scholarship", annualAmount)
	if err != nil {
		t.Fatalf("CalculateMonthlyStipendForStudent failed: %v", err)
	}

	// Assertions
	if result.BaseStipendAmount != expectedMonthly {
		t.Errorf("Monthly base amount mismatch: expected %.2f, got %.2f", expectedMonthly, result.BaseStipendAmount)
	}

	// Cleanup
	cleanupTestDB(t, []models.DeductionRule{rule})
}

// TestCreateStipendForStudent tests stipend creation
func TestCreateStipendForStudent(t *testing.T) {
	// Setup
	setupTestDBNew(t)
	service := NewStipendService()

	studentID := getTestStudentID()
	amount := 50000.0
	stipendType := "full-scholarship"

	// Create stipend
	stipend, err := service.CreateStipendForStudent(
		studentID,
		stipendType,
		amount,
		"Bank_transfer",
		"JN-001-2024",
		"Test stipend",
	)
	if err != nil {
		t.Fatalf("CreateStipendForStudent failed: %v", err)
	}

	// Assertions
	if stipend.StudentID != studentID {
		t.Errorf("Student ID mismatch: expected %s, got %s", studentID, stipend.StudentID)
	}

	if stipend.Amount != amount {
		t.Errorf("Amount mismatch: expected %.2f, got %.2f", amount, stipend.Amount)
	}

	if stipend.StipendType != stipendType {
		t.Errorf("Stipend type mismatch: expected %s, got %s", stipendType, stipend.StipendType)
	}

	if stipend.PaymentStatus != "Pending" {
		t.Errorf("Payment status should be Pending, got %s", stipend.PaymentStatus)
	}

	// Cleanup
	database.DB.Delete(&models.Stipend{}, "id = ?", stipend.ID)
}

// TestStipendValidation tests input validation for stipend creation
func TestStipendValidation(t *testing.T) {
	setupTestDBNew(t)
	service := NewStipendService()

	testCases := []struct {
		name        string
		studentID   uuid.UUID
		stipendType string
		amount      float64
		shouldFail  bool
	}{
		{
			name:        "Valid stipend",
			studentID:   uuid.New(),
			stipendType: "full-scholarship",
			amount:      50000.0,
			shouldFail:  false,
		},
		{
			name:        "Invalid student ID",
			studentID:   uuid.Nil,
			stipendType: "full-scholarship",
			amount:      50000.0,
			shouldFail:  true,
		},
		{
			name:        "Invalid stipend type",
			studentID:   uuid.New(),
			stipendType: "invalid",
			amount:      50000.0,
			shouldFail:  true,
		},
		{
			name:        "Negative amount",
			studentID:   uuid.New(),
			stipendType: "full-scholarship",
			amount:      -1000.0,
			shouldFail:  true,
		},
		{
			name:        "Excessive amount",
			studentID:   uuid.New(),
			stipendType: "full-scholarship",
			amount:      15000000.0,
			shouldFail:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateStipendForStudent(
				tc.studentID,
				tc.stipendType,
				tc.amount,
				"Bank_transfer",
				"JN-"+uuid.New().String(),
				"",
			)

			if tc.shouldFail && err == nil {
				t.Error("Expected error but got nil")
			}

			if !tc.shouldFail && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper functions for testing







func cleanupTestDB(t *testing.T, rules []models.DeductionRule) {
	for _, rule := range rules {
		if err := database.DB.Delete(&rule).Error; err != nil && err != gorm.ErrRecordNotFound {
			t.Logf("Failed to cleanup rule %s: %v", rule.ID, err)
		}
	}
}
