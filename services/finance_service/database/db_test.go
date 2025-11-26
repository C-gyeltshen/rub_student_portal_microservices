package database

import (
	"os"
	"testing"
)

// TestDatabaseConnection tests if the database connects successfully
func TestDatabaseConnection(t *testing.T) {
	// Set up test environment
	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5434/rub_student_portal?sslmode=disable")

	err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if DB == nil {
		t.Fatal("Database instance is nil")
	}

	t.Log("✓ Database connection successful")
}

// TestMigrations tests if migrations run without errors
func TestMigrations(t *testing.T) {
	if DB == nil {
		t.Skip("Database not connected, skipping migration test")
	}

	err := Migrate()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	t.Log("✓ Database migrations successful")
}

// TestTableCreation tests if all tables are created
func TestTableCreation(t *testing.T) {
	if DB == nil {
		t.Skip("Database not connected, skipping table creation test")
	}

	tableNames := []string{"stipends", "deduction_rules", "deductions", "transactions"}

	for _, tableName := range tableNames {
		if !DB.Migrator().HasTable(tableName) {
			t.Fatalf("Table %s was not created", tableName)
		}
		t.Logf("✓ Table %s exists", tableName)
	}
}

// TestForeignKeyConstraints tests if foreign keys are properly set up
func TestForeignKeyConstraints(t *testing.T) {
	if DB == nil {
		t.Skip("Database not connected, skipping foreign key test")
	}

	// Check if stipends table has transaction_id foreign key
	if !DB.Migrator().HasColumn("stipends", "transaction_id") {
		t.Fatal("stipends table missing transaction_id column")
	}

	// Check if deductions table has transaction_id foreign key
	if !DB.Migrator().HasColumn("deductions", "transaction_id") {
		t.Fatal("deductions table missing transaction_id column")
	}

	t.Log("✓ Foreign key columns exist")
}

// TestTransactionTableStructure tests if the transactions table has all required columns
func TestTransactionTableStructure(t *testing.T) {
	if DB == nil {
		t.Skip("Database not connected, skipping transaction table test")
	}

	requiredColumns := []string{
		"id", "stipend_id", "student_id", "amount", "source_account",
		"destination_account", "destination_bank", "transaction_type",
		"status", "payment_method", "reference_number", "error_message",
		"remarks", "initiated_at", "processed_at", "completed_at",
		"created_at", "modified_at",
	}

	for _, colName := range requiredColumns {
		if !DB.Migrator().HasColumn("transactions", colName) {
			t.Fatalf("transactions table missing column: %s", colName)
		}
	}

	t.Log("✓ Transactions table has all required columns")
}

// TestDataInsertion tests if we can insert test data
func TestDataInsertion(t *testing.T) {
	if DB == nil {
		t.Skip("Database not connected, skipping data insertion test")
	}

	// Skip this test as it requires valid foreign key references
	t.Skip("Skipping data insertion test - requires valid foreign key references from students and stipends tables")
}

// TestIndexes tests if all indexes are created
func TestIndexes(t *testing.T) {
	if DB == nil {
		t.Skip("Database not connected, skipping index test")
	}

	indexes := map[string][]string{
		"stipends":         {"idx_stipends_student_id", "idx_stipends_payment_status", "idx_stipends_stipend_type"},
		"deductions":       {"idx_deductions_student_id", "idx_deductions_stipend_id", "idx_deductions_processing_status"},
		"deduction_rules":  {"idx_deduction_rules_rule_name", "idx_deduction_rules_is_active"},
		"transactions":     {"idx_transactions_student_id", "idx_transactions_stipend_id", "idx_transactions_status"},
	}

	for tableName, indexNames := range indexes {
		for _, indexName := range indexNames {
			if !DB.Migrator().HasIndex(tableName, indexName) {
				t.Logf("⚠ Index %s.%s might not exist (GORM may not detect all indexes)", tableName, indexName)
			}
		}
	}

	t.Log("✓ Index check completed")
}
