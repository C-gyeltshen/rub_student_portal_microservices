# Finance Service - Test Documentation

## Overview

This document provides comprehensive documentation for the gRPC integration tests in the Finance Service microservice. The Finance Service handles stipend calculations, deduction management, and financial operations for the RUB Student Portal.

## Test Architecture

### Test Structure

The test suite is organized as follows:

```
services/finance_service/
├── internal/grpc/
│   ├── grpc_test_setup.go          # Common test setup and utilities
│   ├── deduction_server_test.go    # Deduction service gRPC tests
│   ├── stipend_server_test.go      # Stipend service gRPC tests
│   ├── deduction_server.go         # gRPC server implementation
│   └── stipend_server.go           # gRPC server implementation
└── services/
    ├── deduction_service.go        # Deduction business logic
    ├── stipend_service.go          # Stipend business logic
    ├── test_helpers.go             # Helper functions for service tests
    └── test_main.go                # Service-level test setup
```

### Testing Approach

The Finance Service uses **gRPC integration tests** that:

- Connect to a running gRPC server (not unit tests)
- Test the complete request/response cycle
- Validate data persistence in the database
- Ensure proper error handling

## Running the Tests

### Prerequisites

1. **Database**: PostgreSQL running with the `rub_student_portal` database

   ```bash
   docker-compose up postgres
   ```

2. **Environment Variable**: Set the database connection string

   ```bash
   export DATABASE_URL="postgresql://postgres:postgres@localhost:5434/rub_student_portal?sslmode=disable"
   ```

3. **gRPC Server**: The server must be running
   ```bash
   cd services/finance_service
   ./finance_service  # Or: go run main.go
   ```

### Running All Tests

```bash
cd services/finance_service
export DATABASE_URL="postgresql://postgres:postgres@localhost:5434/rub_student_portal?sslmode=disable"
go test ./internal/grpc -v
```

### Running Specific Test

```bash
# Run a single test
go test ./internal/grpc -v -run "TestDeductionService_ListDeductionRules"

# Run all deduction tests
go test ./internal/grpc -v -run "TestDeductionService"

# Run all stipend tests
go test ./internal/grpc -v -run "TestStipendService"
```

### With Coverage Report

```bash
go test ./internal/grpc -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Categories

### 1. Deduction Service Tests

#### TestDeductionService_CreateDeductionRule

- **Purpose**: Verify creation of deduction rules
- **What it tests**:
  - Creating a new deduction rule with all parameters
  - Rule is marked as active upon creation
  - Rule details match the request
- **Data used**: Test rule with hostel fee deduction
- **Expected outcome**: Rule created successfully with ID

#### TestDeductionService_GetDeductionRule

- **Purpose**: Verify retrieval of a single deduction rule
- **What it tests**:
  - Creating a rule and retrieving it by ID
  - Retrieved data matches created data
  - Correct rule name and deduction type
- **Data used**: Electricity bill deduction rule
- **Expected outcome**: Rule retrieved with matching details

#### TestDeductionService_ListDeductionRules

- **Purpose**: Verify pagination and filtering of deduction rules
- **What it tests**:
  - List returns paginated results
  - Total count is accurate
  - Rules are filtered by active status
  - Respects limit and offset parameters
- **Data used**: Multiple pre-existing rules
- **Expected outcome**: Paginated list with correct total

#### TestDeductionService_CreateDeduction

- **Purpose**: Verify creation of deductions for a stipend
- **What it tests**:
  - Creating a stipend first
  - Creating a deduction rule
  - Creating a deduction linking student, rule, and stipend
  - Amount is properly stored
- **Data used**: Test student, full-scholarship stipend, mess fee deduction
- **Expected outcome**: Deduction created with correct amount

#### TestDeductionService_GetDeduction

- **Purpose**: Verify retrieval of a single deduction
- **What it tests**:
  - Creating deduction infrastructure (stipend, rule, deduction)
  - Retrieving deduction by ID
  - Retrieved details match created data
- **Data used**: Water bill deduction
- **Expected outcome**: Deduction retrieved with matching details

#### TestDeductionService_GetStipendDeductions

- **Purpose**: Verify retrieving all deductions for a specific stipend
- **What it tests**:
  - Creating multiple deductions for one stipend
  - Listing deductions with pagination
  - Correct total count of deductions
  - Deductions ordered by date (descending)
- **Data used**: Multiple deductions for same stipend
- **Expected outcome**: List of deductions for stipend with correct total

#### TestDeductionService_ApplyDeductions

- **Purpose**: Verify the bulk application of deductions
- **What it tests**:
  - Applying multiple deduction rules to a single stipend
  - Each deduction is created
  - Total deduction amount is calculated correctly
  - Deductions respect min/max constraints
- **Data used**: Stipend with multiple applicable rules
- **Expected outcome**: All deductions applied with correct amounts

### 2. Stipend Service Tests

#### TestStipendService_CalculateStipendWithDeductions

- **Purpose**: Verify stipend calculation with applicable deductions
- **What it tests**:
  - Base stipend amount is correct
  - All applicable deductions are identified
  - Deductions are applied in priority order
  - Net stipend (base - deductions) is calculated correctly
  - Deduction details are returned
- **Data used**: Full-scholarship stipend type
- **Expected outcome**: Calculation result with base, deductions, and net amount

#### TestStipendService_CalculateMonthlyStipend

- **Purpose**: Verify monthly stipend calculation
- **What it tests**:
  - Annual stipend divided by 12 (approximately)
  - Monthly deductions applied
  - Net monthly amount calculated
  - Monthly-applicable rules are used
- **Data used**: Full-scholarship annual amount divided monthly
- **Expected outcome**: Monthly stipend with monthly deductions

#### TestStipendService_CalculateAnnualStipend

- **Purpose**: Verify annual stipend calculation
- **What it tests**:
  - Annual base amount used
  - Annual and mandatory deductions applied
  - Net annual amount calculated
  - Correct deduction rules applied
- **Data used**: Full-scholarship annual stipend
- **Expected outcome**: Annual stipend with all deductions

#### TestStipendService_CreateStipend

- **Purpose**: Verify creation of stipend records
- **What it tests**:
  - Stipend is created with correct student
  - Amount is stored correctly
  - Payment status is set to "Pending"
  - Stipend type is stored
  - Journal number is unique
- **Data used**: Test student, 100,000 amount, bank transfer payment
- **Expected outcome**: Stipend created with all details

#### TestStipendService_GetStipend

- **Purpose**: Verify retrieval of a stipend by ID
- **What it tests**:
  - Creating a stipend
  - Retrieving it by ID
  - Retrieved details match created data
  - All fields are correctly returned
- **Data used**: Previously created stipend
- **Expected outcome**: Stipend retrieved with matching details

#### TestStipendService_UpdateStipendPaymentStatus

- **Purpose**: Verify updating stipend payment status
- **What it tests**:
  - Creating a stipend with "Pending" status
  - Updating status to "Processed"
  - Updated status is persisted
  - Original other details unchanged
- **Data used**: Stipend status transition
- **Expected outcome**: Status successfully updated

## Test Setup and Utilities

### grpc_test_setup.go

This file contains the common setup and utility functions for all gRPC tests.

#### Key Functions

**TestMain(m \*testing.M)**

- Executed before all tests in the package
- Sets up database connection
- Runs migrations
- Initializes finance database (indexes, constraints, seeds)
- Cleans up old test data
- **Note**: Prints initialization status to stdout

**GetTestStudentID() string**

- Lazy-initializes a test student in the database
- Returns existing student ID if available
- Creates new test student if needed
- Ensures foreign key constraints are satisfied
- Called by each test that needs a student

**initializeTestStudent() error**

- Attempts to find existing test student
- Creates new student if none exists
- Generates unique email and card number
- Retries up to 5 times on conflict
- Returns error if all attempts fail

**cleanupOldTestData() error**

- Removes test deduction rules older than 1 hour
- Prevents database bloat from repeated test runs
- Uses pattern matching on rule names

**GetUniqueTestName(baseName string) string**

- Generates unique names/IDs using Unix nanosecond timestamps
- Prevents constraint violations from duplicate data
- Format: `{baseName}_{timestamp}`

## Data Management in Tests

### Test Data Lifecycle

```
Test Start
    ↓
Initialize Database Connection
    ↓
Create/Find Test Student
    ↓
Run Test Operations
    ├─ Create Stipend
    ├─ Create Deduction Rule
    ├─ Create Deduction
    └─ Verify Results
    ↓
Test End
```

### Key Test Data

**Test Student**

- Auto-created with unique email: `test-student-{timestamp}@rub.edu.bt`
- Unique card number: `TEST{timestamp}`
- Phone number: `+97517123456`

**Test Stipend**

- Amount: 100,000 (or calculated)
- Type: "full-scholarship" or "self-funded"
- Status: "Pending"
- Journal Number: `JN-{type}-{timestamp}` (must be unique per run)

**Test Deduction Rules** (created during each test)

- Names: `Test {Type} {timestamp}` (unique per run)
- Types: hostel, electricity, mess, water, library
- Amounts: Vary by type (1,000 - 5,000)
- Priority: Varies for testing order

### Unique Data Strategy

To avoid constraint violations across test runs:

1. **Journal Numbers**: Use `GetUniqueTestName()` function

   ```go
   JournalNumber: "JN-DED-" + GetUniqueTestName("001")
   ```

2. **Rule Names**: Include timestamp

   ```go
   RuleName: "Test Mess Fee " + GetUniqueTestName("")
   ```

3. **Student Data**: Auto-generated with timestamp
   ```go
   email: fmt.Sprintf("test-student-%d@rub.edu.bt", time.Now().UnixNano())
   ```

## Common Issues and Solutions

### Issue 1: "Cannot connect to gRPC server"

**Cause**: Server not running on port 50051
**Solution**:

```bash
# Start the server
cd services/finance_service
./finance_service
# Or in background: ./finance_service &
```

### Issue 2: "Table not set" GORM error

**Cause**: Missing `.Model()` in GORM query chain
**Solution**: Ensure `.Model(&models.ModelName{})` is called before `.Count()`

```go
// ❌ Wrong
query.Count(&total)

// ✅ Correct
query.Model(&models.Deduction{}).Count(&total)
```

### Issue 3: "Duplicate key value violates constraint"

**Cause**: Test data reused from previous runs
**Solution**: Use `GetUniqueTestName()` for all unique fields:

- Journal numbers
- Rule names
- Student emails

### Issue 4: "Foreign key constraint violation"

**Cause**: Referenced student/rule doesn't exist
**Solution**: Ensure `GetTestStudentID()` is called to initialize student before creating stipends

### Issue 5: "Database connection refused"

**Cause**: PostgreSQL not running or DATABASE_URL incorrect
**Solution**:

```bash
# Check Docker is running
docker ps | grep postgres

# Set correct DATABASE_URL
export DATABASE_URL="postgresql://postgres:postgres@localhost:5434/rub_student_portal?sslmode=disable"
```

## Test Metrics

### Current Test Coverage

| Component         | Tests  | Status           |
| ----------------- | ------ | ---------------- |
| Deduction Service | 7      | ✅ All Pass      |
| Stipend Service   | 6      | ✅ All Pass      |
| **Total**         | **13** | **✅ 100% Pass** |

### Performance Metrics (from recent run)

```
TestDeductionService_CreateDeductionRule        0.01s ✅
TestDeductionService_GetDeductionRule           0.01s ✅
TestDeductionService_ListDeductionRules         0.00s ✅
TestDeductionService_CreateDeduction            0.02s ✅
TestDeductionService_GetDeduction               0.01s ✅
TestDeductionService_GetStipendDeductions       0.01s ✅
TestDeductionService_ApplyDeductions            0.23s ✅
TestStipendService_CalculateStipendWithDeductions 0.00s ✅
TestStipendService_CalculateMonthlyStipend      0.00s ✅
TestStipendService_CalculateAnnualStipend       0.01s ✅
TestStipendService_CreateStipend                0.01s ✅
TestStipendService_GetStipend                   0.01s ✅
TestStipendService_UpdateStipendPaymentStatus   0.01s ✅
─────────────────────────────────────────────────────
Total Duration                                  0.35s
```

## Best Practices for Writing New Tests

### 1. Use Helper Functions

```go
// ✅ Good
studentID := GetTestStudentID()
journalNum := "JN-TEST-" + GetUniqueTestName("")

// ❌ Avoid
studentID := "12345678-1234-1234-1234-123456789012"  // Hard-coded UUID
journalNum := "JN-TEST-001"  // Will conflict on re-run
```

### 2. Set Up Context Properly

```go
ctx := context.Background()
conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    t.Skipf("Could not connect to gRPC server: %v", err)
}
defer conn.Close()
```

### 3. Follow Test Structure

```go
func TestFeatureName(t *testing.T) {
    // 1. Setup
    conn := // ... connect to server
    client := // ... create client

    // 2. Create prerequisites
    student := GetTestStudentID()

    // 3. Execute test
    resp, err := client.Operation(ctx, &Request{})

    // 4. Verify results
    if err != nil {
        t.Fatalf("Operation failed: %v", err)
    }
    if resp.Field != expectedValue {
        t.Errorf("Expected %v, got %v", expectedValue, resp.Field)
    }
}
```

### 4. Handle Unique Constraints

```go
// Create unique rule name for each test run
ruleReq := &pb.CreateDeductionRuleRequest{
    RuleName: "Test Rule " + GetUniqueTestName(""),
    // ... other fields
}
```

## Continuous Integration

### GitHub Actions Configuration Example

```yaml
name: Finance Service Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: rub_student_portal
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: "1.21"

      - name: Start gRPC Server
        working-directory: services/finance_service
        env:
          DATABASE_URL: postgresql://postgres:postgres@localhost:5432/rub_student_portal?sslmode=disable
        run: ./finance_service &

      - name: Run Tests
        working-directory: services/finance_service
        env:
          DATABASE_URL: postgresql://postgres:postgres@localhost:5432/rub_student_portal?sslmode=disable
        run: go test ./internal/grpc -v
```

## Troubleshooting Guide

### Test Hangs

**Symptom**: Test takes >5 seconds or never completes

**Possible causes**:

1. Database connection timeout
2. gRPC server not responding
3. Deadlock in test code

**Solution**:

```bash
# Run with timeout
go test ./internal/grpc -v -timeout 30s

# See verbose logs
go test ./internal/grpc -v -run TestName
```

### Flaky Tests

**Symptom**: Tests pass sometimes, fail other times

**Common causes**:

1. Database constraints violated
2. Race conditions in data setup
3. Timing issues

**Solution**:

- Ensure unique test data using timestamps
- Use database transactions for test isolation
- Add small delays before assertions if needed

### Database Issues

**Symptom**: "connection refused" or "table not found"

**Debug**:

```bash
# Check if database is accessible
docker-compose exec postgres psql -U postgres -c "\l"

# Verify migrations ran
docker-compose exec postgres psql -U postgres rub_student_portal -c "\dt"

# Check test database state
docker-compose exec postgres psql -U postgres rub_student_portal -c "SELECT COUNT(*) FROM students"
```

## Further Reading

- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Go Testing Package](https://pkg.go.dev/testing)
- [GORM Querying](https://gorm.io/docs/query.html)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Contributing Tests

When adding new tests:

1. Follow the existing naming convention: `Test{Service}_{Feature}`
2. Use `GetUniqueTestName()` for all unique data
3. Clean up resources in defer statements
4. Add comments explaining what is being tested
5. Update this documentation with new test descriptions
6. Ensure all tests pass locally before pushing
7. Run full test suite: `go test ./... -v`

---

**Last Updated**: November 25, 2025  
**Test Status**: ✅ All 13 tests passing  
**Coverage**: Deduction & Stipend Services
