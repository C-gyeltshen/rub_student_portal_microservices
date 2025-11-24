package database

import (
	"log"

	"finance_service/models"

	"github.com/google/uuid"
)

// InitializeFinanceDatabase creates indexes and ensures constraints for finance service
func InitializeFinanceDatabase() error {
	log.Println("Initializing finance service database...")

	// Create indexes for performance optimization
	if err := createIndexes(); err != nil {
		log.Printf("Error creating indexes: %v", err)
		return err
	}

	// Ensure foreign key constraints
	if err := ensureForeignKeyConstraints(); err != nil {
		log.Printf("Error ensuring foreign key constraints: %v", err)
		// Don't fail completely - constraints may already exist
	}

	// Seed default deduction rules
	if err := SeedDeductionRules(); err != nil {
		log.Printf("Error seeding deduction rules: %v", err)
		return err
	}

	log.Println("Finance database initialization completed successfully!")
	return nil
}

// createIndexes creates all necessary indexes for finance service tables
func createIndexes() error {
	log.Println("Creating database indexes...")

	indexes := []struct {
		table  string
		column string
		name   string
	}{
		// Stipends indexes
		{"stipends", "student_id", "idx_stipends_student_id"},
		{"stipends", "payment_status", "idx_stipends_payment_status"},
		{"stipends", "stipend_type", "idx_stipends_stipend_type"},
		{"stipends", "created_at", "idx_stipends_created_at"},

		// Deductions indexes
		{"deductions", "student_id", "idx_deductions_student_id"},
		{"deductions", "stipend_id", "idx_deductions_stipend_id"},
		{"deductions", "processing_status", "idx_deductions_processing_status"},
		{"deductions", "deduction_type", "idx_deductions_deduction_type"},
		{"deductions", "deduction_rule_id", "idx_deductions_deduction_rule_id"},

		// Deduction Rules indexes
		{"deduction_rules", "rule_name", "idx_deduction_rules_rule_name"},
		{"deduction_rules", "is_active", "idx_deduction_rules_is_active"},
		{"deduction_rules", "deduction_type", "idx_deduction_rules_deduction_type"},
	}

	for _, idx := range indexes {
		if err := DB.Migrator().CreateIndex(idx.table, idx.column); err != nil {
			log.Printf("Warning: Could not create index %s on %s.%s: %v", idx.name, idx.table, idx.column, err)
			// Continue - index may already exist
		} else {
			log.Printf("✓ Created index: %s on %s(%s)", idx.name, idx.table, idx.column)
		}
	}

	return nil
}

// ensureForeignKeyConstraints adds foreign key constraints that might be missing
func ensureForeignKeyConstraints() error {
	log.Println("Ensuring foreign key constraints...")

	constraints := []struct {
		table          string
		column         string
		refTable       string
		refColumn      string
		constraintName string
		onDelete       string
	}{
		{
			table:          "stipends",
			column:         "student_id",
			refTable:       "students",
			refColumn:      "id",
			constraintName: "fk_stipends_student_id",
			onDelete:       "CASCADE",
		},
		{
			table:          "stipends",
			column:         "transaction_id",
			refTable:       "transaction",
			refColumn:      "id",
			constraintName: "fk_stipends_transaction_id",
			onDelete:       "SET NULL",
		},
		{
			table:          "deductions",
			column:         "student_id",
			refTable:       "students",
			refColumn:      "id",
			constraintName: "fk_deductions_student_id",
			onDelete:       "CASCADE",
		},
		{
			table:          "deductions",
			column:         "stipend_id",
			refTable:       "stipends",
			refColumn:      "id",
			constraintName: "fk_deductions_stipend_id",
			onDelete:       "CASCADE",
		},
		{
			table:          "deductions",
			column:         "deduction_rule_id",
			refTable:       "deduction_rules",
			refColumn:      "id",
			constraintName: "fk_deductions_deduction_rule_id",
			onDelete:       "RESTRICT",
		},
		{
			table:          "deductions",
			column:         "approved_by",
			refTable:       "users",
			refColumn:      "id",
			constraintName: "fk_deductions_approved_by",
			onDelete:       "SET NULL",
		},
		{
			table:          "deductions",
			column:         "transaction_id",
			refTable:       "transaction",
			refColumn:      "id",
			constraintName: "fk_deductions_transaction_id",
			onDelete:       "SET NULL",
		},
		{
			table:          "deduction_rules",
			column:         "created_by",
			refTable:       "users",
			refColumn:      "id",
			constraintName: "fk_deduction_rules_created_by",
			onDelete:       "SET NULL",
		},
		{
			table:          "deduction_rules",
			column:         "modified_by",
			refTable:       "users",
			refColumn:      "id",
			constraintName: "fk_deduction_rules_modified_by",
			onDelete:       "SET NULL",
		},
	}

	for _, constraint := range constraints {
		// Check if constraint already exists
		if DB.Migrator().HasConstraint(constraint.table, constraint.constraintName) {
			log.Printf("✓ Constraint already exists: %s", constraint.constraintName)
			continue
		}

		// Add the foreign key constraint
		if err := DB.Migrator().CreateConstraint(constraint.table, constraint.column); err != nil {
			log.Printf("Warning: Could not create constraint %s: %v", constraint.constraintName, err)
			// Continue - constraint may already exist or be unnecessary
		} else {
			log.Printf("✓ Created constraint: %s", constraint.constraintName)
		}
	}

	return nil
}

// SeedDeductionRules populates the database with initial deduction rules
func SeedDeductionRules() error {
	log.Println("Seeding deduction rules...")

	rules := []models.DeductionRule{
		// Rules applicable to all students
		{
			ID:                        uuid.New(),
			RuleName:                  "Hostel Fee",
			DeductionType:             "hostel",
			Description:               "Monthly on-campus hostel accommodation charges",
			BaseAmount:                3000.00,
			MaxDeductionAmount:        3500.00,
			MinDeductionAmount:        2500.00,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			AppliesMonthly:            true,
			AppliesAnnually:           false,
			IsOptional:                false,
			Priority:                  100,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "Electricity Bill",
			DeductionType:             "electricity",
			Description:               "Monthly electricity charges for hostel rooms",
			BaseAmount:                500.00,
			MaxDeductionAmount:        800.00,
			MinDeductionAmount:        300.00,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			AppliesMonthly:            true,
			AppliesAnnually:           false,
			IsOptional:                false,
			Priority:                  80,
		},

		// Rules for self-funded students only
		{
			ID:                        uuid.New(),
			RuleName:                  "Mess Fee",
			DeductionType:             "mess_fees",
			Description:               "Monthly dining facility and meal charges",
			BaseAmount:                2000.00,
			MaxDeductionAmount:        2500.00,
			MinDeductionAmount:        1500.00,
			IsApplicableToFullScholar: false,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			AppliesMonthly:            true,
			AppliesAnnually:           false,
			IsOptional:                false,
			Priority:                  90,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "Water Bill",
			DeductionType:             "water",
			Description:               "Monthly water supply charges",
			BaseAmount:                300.00,
			MaxDeductionAmount:        500.00,
			MinDeductionAmount:        200.00,
			IsApplicableToFullScholar: false,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			AppliesMonthly:            true,
			AppliesAnnually:           false,
			IsOptional:                false,
			Priority:                  70,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "Library Fine",
			DeductionType:             "library",
			Description:               "Library facility and book damage charges",
			BaseAmount:                200.00,
			MaxDeductionAmount:        500.00,
			MinDeductionAmount:        0.00,
			IsApplicableToFullScholar: false,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			AppliesMonthly:            false,
			AppliesAnnually:           true,
			IsOptional:                true,
			Priority:                  30,
		},

		// Optional charges
		{
			ID:                        uuid.New(),
			RuleName:                  "Sports Activity Fee",
			DeductionType:             "sports",
			Description:               "Optional sports and recreational facilities fee",
			BaseAmount:                500.00,
			MaxDeductionAmount:        1000.00,
			MinDeductionAmount:        0.00,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			AppliesMonthly:            true,
			AppliesAnnually:           false,
			IsOptional:                true,
			Priority:                  20,
		},
		{
			ID:                        uuid.New(),
			RuleName:                  "University Fund Contribution",
			DeductionType:             "university_fund",
			Description:               "Contribution to university development fund",
			BaseAmount:                1000.00,
			MaxDeductionAmount:        2000.00,
			MinDeductionAmount:        500.00,
			IsApplicableToFullScholar: true,
			IsApplicableToSelfFunded:  true,
			IsActive:                  true,
			AppliesMonthly:            false,
			AppliesAnnually:           true,
			IsOptional:                true,
			Priority:                  10,
		},
	}

	for _, rule := range rules {
		// Check if rule already exists
		var count int64
		if err := DB.Model(&models.DeductionRule{}).Where("rule_name = ?", rule.RuleName).Count(&count).Error; err != nil {
			log.Printf("Error checking existing rule %s: %v", rule.RuleName, err)
			continue
		}

		if count > 0 {
			log.Printf("Rule already exists: %s", rule.RuleName)
			continue
		}

		// Create the rule
		if err := DB.Create(&rule).Error; err != nil {
			log.Printf("Error creating deduction rule %s: %v", rule.RuleName, err)
			continue
		}

		log.Printf("✓ Created deduction rule: %s (Priority: %d, Applicable to Full-Scholarship: %v, Self-Funded: %v)",
			rule.RuleName, rule.Priority, rule.IsApplicableToFullScholar, rule.IsApplicableToSelfFunded)
	}

	log.Println("Deduction rules seeding completed!")
	return nil
}
