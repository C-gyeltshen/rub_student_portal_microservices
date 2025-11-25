package services

import (
	"os"
	"testing"

	"finance_service/database"

	"github.com/google/uuid"
)

// TestCreateDeductionRule tests rule creation
func TestCreateDeductionRule(t *testing.T) {
	// Ensure database is initialized
	if database.DB == nil {
		if os.Getenv("DATABASE_URL") == "" {
			os.Setenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/rub_student_portal?sslmode=disable")
		}
		if err := database.Connect(); err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		if err := database.Migrate(); err != nil {
			t.Fatalf("Failed to migrate database: %v", err)
		}
	}

	service := NewDeductionRuleService()

	input := &CreateDeductionRuleInput{
		RuleName:                  "Hostel Fee Rule",
		DeductionType:             "hostel",
		Description:               "Monthly hostel fee deduction",
		BaseAmount:                5000.00,
		MaxDeductionAmount:        5000.00,
		MinDeductionAmount:        0.00,
		IsApplicableToFullScholar: false,
		IsApplicableToSelfFunded:  true,
		AppliesMonthly:            true,
		AppliesAnnually:           false,
		IsOptional:                false,
		Priority:                  1,
	}

	rule, err := service.CreateRule(input)
	if err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	if rule.RuleName != input.RuleName {
		t.Errorf("Expected rule name %s, got %s", input.RuleName, rule.RuleName)
	}

	if rule.BaseAmount != input.BaseAmount {
		t.Errorf("Expected base amount %.2f, got %.2f", input.BaseAmount, rule.BaseAmount)
	}

	if !rule.IsActive {
		t.Errorf("Expected rule to be active")
	}
}

// TestCreateDuplicateRuleName tests that duplicate rule names are rejected
func TestCreateDuplicateRuleName(t *testing.T) {
	service := NewDeductionRuleService()

	input := &CreateDeductionRuleInput{
		RuleName:                  "Unique Hostel Fee",
		DeductionType:             "hostel",
		Description:               "Monthly hostel fee",
		BaseAmount:                5000.00,
		MaxDeductionAmount:        5000.00,
		MinDeductionAmount:        0.00,
		IsApplicableToFullScholar: false,
		IsApplicableToSelfFunded:  true,
		AppliesMonthly:            true,
		AppliesAnnually:           false,
		IsOptional:                false,
		Priority:                  1,
	}

	// Create first rule
	_, err := service.CreateRule(input)
	if err != nil {
		t.Fatalf("Failed to create first rule: %v", err)
	}

	// Try to create duplicate
	_, err = service.CreateRule(input)
	if err == nil {
		t.Errorf("Expected error for duplicate rule name, got nil")
	}
}

// TestValidationErrors tests input validation
func TestValidationErrors(t *testing.T) {
	service := NewDeductionRuleService()

	tests := []struct {
		name    string
		input   *CreateDeductionRuleInput
		wantErr bool
	}{
		{
			name: "Missing rule name",
			input: &CreateDeductionRuleInput{
				RuleName:       "",
				DeductionType:  "hostel",
				BaseAmount:     1000,
				AppliesMonthly: true,
			},
			wantErr: true,
		},
		{
			name: "Missing deduction type",
			input: &CreateDeductionRuleInput{
				RuleName:       "Test Rule",
				DeductionType:  "",
				BaseAmount:     1000,
				AppliesMonthly: true,
			},
			wantErr: true,
		},
		{
			name: "Negative base amount",
			input: &CreateDeductionRuleInput{
				RuleName:       "Test Rule",
				DeductionType:  "hostel",
				BaseAmount:     -1000,
				AppliesMonthly: true,
			},
			wantErr: true,
		},
		{
			name: "Max less than min deduction",
			input: &CreateDeductionRuleInput{
				RuleName:           "Test Rule",
				DeductionType:      "hostel",
				BaseAmount:         1000,
				MaxDeductionAmount: 500,
				MinDeductionAmount: 1000,
				AppliesMonthly:     true,
			},
			wantErr: true,
		},
		{
			name: "No frequency specified",
			input: &CreateDeductionRuleInput{
				RuleName:       "Test Rule",
				DeductionType:  "hostel",
				BaseAmount:     1000,
				AppliesMonthly: false,
				AppliesAnnually: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateRule(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUpdateDeductionRule tests rule updates
func TestUpdateDeductionRule(t *testing.T) {
	service := NewDeductionRuleService()

	// Create a rule
	input := &CreateDeductionRuleInput{
		RuleName:                  "Original Name",
		DeductionType:             "electricity",
		Description:               "Original description",
		BaseAmount:                1000,
		MaxDeductionAmount:        1000,
		MinDeductionAmount:        0,
		IsApplicableToFullScholar: false,
		IsApplicableToSelfFunded:  true,
		AppliesMonthly:            true,
		AppliesAnnually:           false,
		IsOptional:                false,
		Priority:                  1,
	}

	rule, err := service.CreateRule(input)
	if err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	// Update the rule
	newName := "Updated Name"
	newAmount := 2000.0
	updateInput := &UpdateDeductionRuleInput{
		RuleName:   &newName,
		BaseAmount: &newAmount,
	}

	updatedRule, err := service.UpdateRule(rule.ID, updateInput)
	if err != nil {
		t.Fatalf("Failed to update rule: %v", err)
	}

	if updatedRule.RuleName != newName {
		t.Errorf("Expected rule name %s, got %s", newName, updatedRule.RuleName)
	}

	if updatedRule.BaseAmount != newAmount {
		t.Errorf("Expected base amount %.2f, got %.2f", newAmount, updatedRule.BaseAmount)
	}
}

// TestDeleteDeductionRule tests rule deletion (soft delete)
func TestDeleteDeductionRule(t *testing.T) {
	service := NewDeductionRuleService()

	// Create a rule
	input := &CreateDeductionRuleInput{
		RuleName:       "Rule to Delete",
		DeductionType:  "mess",
		BaseAmount:     3000,
		AppliesMonthly: true,
	}

	rule, err := service.CreateRule(input)
	if err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	// Verify rule exists and is active
	retrievedRule, err := service.GetRuleByID(rule.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve rule: %v", err)
	}
	if !retrievedRule.IsActive {
		t.Errorf("Expected rule to be active initially")
	}

	// Delete the rule
	err = service.DeleteRule(rule.ID)
	if err != nil {
		t.Fatalf("Failed to delete rule: %v", err)
	}

	// Verify rule is now inactive
	retrievedRule, err = service.GetRuleByID(rule.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve rule after delete: %v", err)
	}
	if retrievedRule.IsActive {
		t.Errorf("Expected rule to be inactive after deletion")
	}
}

// TestGetRuleByID tests retrieving a rule by ID
func TestGetRuleByID(t *testing.T) {
	service := NewDeductionRuleService()

	// Create a rule
	input := &CreateDeductionRuleInput{
		RuleName:       "Get By ID Test",
		DeductionType:  "test",
		BaseAmount:     1500,
		AppliesMonthly: true,
	}

	rule, err := service.CreateRule(input)
	if err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	// Retrieve by ID
	retrieved, err := service.GetRuleByID(rule.ID)
	if err != nil {
		t.Fatalf("Failed to get rule by ID: %v", err)
	}

	if retrieved.ID != rule.ID {
		t.Errorf("Expected ID %s, got %s", rule.ID, retrieved.ID)
	}
	if retrieved.RuleName != rule.RuleName {
		t.Errorf("Expected name %s, got %s", rule.RuleName, retrieved.RuleName)
	}

	// Try to get non-existent rule
	_, err = service.GetRuleByID(uuid.New())
	if err == nil {
		t.Errorf("Expected error for non-existent rule, got nil")
	}
}

// TestListActiveRules tests listing active rules
func TestListActiveRules(t *testing.T) {
	service := NewDeductionRuleService()

	// Create multiple rules
	for i := 0; i < 3; i++ {
		input := &CreateDeductionRuleInput{
			RuleName:       "Active Rule " + string(rune('A'+i)),
			DeductionType:  "hostel",
			BaseAmount:     float64(1000 + i*100),
			AppliesMonthly: true,
			Priority:       i,
		}
		_, err := service.CreateRule(input)
		if err != nil {
			t.Fatalf("Failed to create rule: %v", err)
		}
	}

	// List active rules
	rules, total, err := service.ListActiveRules(10, 0)
	if err != nil {
		t.Fatalf("Failed to list active rules: %v", err)
	}

	if len(rules) < 3 {
		t.Errorf("Expected at least 3 active rules, got %d", len(rules))
	}

	if total < 3 {
		t.Errorf("Expected total >= 3, got %d", total)
	}

	// Verify all returned rules are active
	for _, rule := range rules {
		if !rule.IsActive {
			t.Errorf("Expected all listed rules to be active")
		}
	}
}

// TestListRulesByType tests filtering rules by type
func TestListRulesByType(t *testing.T) {
	service := NewDeductionRuleService()

	// Create rules of different types
	types := []string{"hostel", "electricity", "mess"}
	for _, ruleType := range types {
		input := &CreateDeductionRuleInput{
			RuleName:       ruleType + " rule",
			DeductionType:  ruleType,
			BaseAmount:     1000,
			AppliesMonthly: true,
		}
		_, err := service.CreateRule(input)
		if err != nil {
			t.Fatalf("Failed to create rule: %v", err)
		}
	}

	// List rules by type
	rules, _, err := service.ListRulesByType("hostel", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list rules by type: %v", err)
	}

	// Verify all returned rules are of the requested type
	for _, rule := range rules {
		if rule.DeductionType != "hostel" {
			t.Errorf("Expected all rules to be type 'hostel', got %s", rule.DeductionType)
		}
	}
}

// TestGetApplicableRules tests getting rules applicable to a specific student type
func TestGetApplicableRules(t *testing.T) {
	service := NewDeductionRuleService()

	// Create rules applicable to different student types
	inputs := []struct {
		name              string
		fullScholar       bool
		selfFunded        bool
	}{
		{"Full Scholar Only", true, false},
		{"Self Funded Only", false, true},
		{"Both Types", true, true},
	}

	for _, input := range inputs {
		rule := &CreateDeductionRuleInput{
			RuleName:                  input.name,
			DeductionType:             "test",
			BaseAmount:                1000,
			IsApplicableToFullScholar: input.fullScholar,
			IsApplicableToSelfFunded:  input.selfFunded,
			AppliesMonthly:            true,
		}
		_, err := service.CreateRule(rule)
		if err != nil {
			t.Fatalf("Failed to create rule: %v", err)
		}
	}

	// Get applicable rules for full scholar
	fullScholarRules, err := service.GetApplicableRules(true)
	if err != nil {
		t.Fatalf("Failed to get applicable rules for full scholar: %v", err)
	}

	if len(fullScholarRules) < 2 {
		t.Errorf("Expected at least 2 rules for full scholar, got %d", len(fullScholarRules))
	}

	// Get applicable rules for self-funded
	selfFundedRules, err := service.GetApplicableRules(false)
	if err != nil {
		t.Fatalf("Failed to get applicable rules for self-funded: %v", err)
	}

	if len(selfFundedRules) < 2 {
		t.Errorf("Expected at least 2 rules for self-funded, got %d", len(selfFundedRules))
	}
}
