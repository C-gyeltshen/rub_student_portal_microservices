# Finance Service - Stipend Calculation & Deduction System

## Overview

The Finance Service manages stipend calculations, deductions, and money transfers for students. It supports:

- **Full-Scholarship and Self-Funded Student Management**: Different deduction rules for each type
- **Configurable Deduction Rules**: Hostel, electricity, mess fees, and more
- **Priority-Based Deduction Ordering**: Apply deductions in configurable priority order
- **Stipend Calculation**: Base amount minus applicable deductions
- **Money Transfer**: Transfer stipends directly to student bank accounts
- **Payment Processing**: Track payment status, transfer history, and transaction details
- **Audit Logging**: Complete history of all stipend, deduction, and transaction records

## Quick Start

See [QUICK_START.md](./QUICK_START.md) for quick setup and usage examples.

See [API_REFERENCE.md](./API_REFERENCE.md) for complete API documentation.

## Key Features (Task 1.2 + Money Transfer Implementation)

### ✅ Stipend Types

- **Full-Scholarship Students**: Limited deductions (configurable per rule)
- **Self-Funded Students**: More deductions apply by default

### ✅ Stipend Calculation

1. Base amount (annual or monthly)
2. Fetch applicable deduction rules for student type
3. Sort by priority (highest first)
4. Apply each deduction in order
5. Cap deductions to remaining stipend
6. Return base, deductions, and net amounts

### ✅ HTTP Endpoints - Stipends & Deductions

- `POST /api/stipends`: Create stipend
- `POST /api/stipends/calculate`: Calculate with deductions
- `POST /api/stipends/calculate/monthly`: Monthly calculation
- `POST /api/stipends/calculate/annual`: Annual calculation
- `GET /api/stipends/{stipendID}`: Get stipend
- `GET /api/students/{studentID}/stipends`: List student stipends
- `PATCH /api/stipends/{stipendID}/payment-status`: Update status
- `GET /api/stipends/{stipendID}/deductions`: Get deductions
- `POST /api/deduction-rules`: Create rule
- `GET /api/deduction-rules`: List rules
- `GET /api/deduction-rules/{ruleID}`: Get rule

### ✅ Money Transfer Endpoints (NEW!)

- `POST /api/transfers/initiate`: Initiate a money transfer to student's bank
- `POST /api/transfers/{transactionID}/process`: Process pending transfer via payment gateway
- `GET /api/transfers/{transactionID}/status`: Check transfer status
- `GET /api/stipends/{stipendID}/transactions`: Get all transactions for a stipend
- `GET /api/students/{studentID}/transactions`: Get all transactions for a student
- `POST /api/transfers/{transactionID}/cancel`: Cancel a pending/processing transfer
- `POST /api/transfers/{transactionID}/retry`: Retry a failed transfer

### ✅ Deduction Management

- Create configurable deduction rules
- Support monthly and annual deductions
- Optional vs. mandatory deductions
- Min/max deduction bounds
- Student-type applicability

## Database Schema

### Tables

#### `stipends`

Stores stipend payment records for students.

**Columns:**

- `id` (UUID): Primary key
- `student_id` (UUID): Foreign key to students table
- `amount` (DECIMAL): Stipend amount
- `stipend_type` (VARCHAR): Type of stipend (full-scholarship, self-funded)
- `payment_date` (TIMESTAMPTZ): When the stipend was paid
- `payment_status` (VARCHAR): Payment status (Pending, Processed, Failed)
- `payment_method` (VARCHAR): Payment method (Bank_transfer, E-payment)
- `journal_number` (TEXT): Unique journal reference number
- `transaction_id` (UUID): Foreign key to transaction table
- `notes` (TEXT): Additional notes
- `created_at` (TIMESTAMPTZ): Record creation timestamp
- `modified_at` (TIMESTAMPTZ): Last modification timestamp

#### `deduction_rules`

Defines configurable deduction rules that can be applied to stipends.

**Columns:**

- `id` (UUID): Primary key
- `rule_name` (VARCHAR): Unique name of the rule
- `deduction_type` (VARCHAR): Type of deduction (hostel, electricity, mess_fees, etc.)
- `description` (TEXT): Rule description
- `base_amount` (DECIMAL): Base amount for this deduction
- `max_deduction_amount` (DECIMAL): Maximum deduction allowed
- `min_deduction_amount` (DECIMAL): Minimum deduction amount
- `is_applicable_to_full_scholar` (BOOLEAN): Whether rule applies to full-scholarship students
- `is_applicable_to_self_funded` (BOOLEAN): Whether rule applies to self-funded students
- `is_active` (BOOLEAN): Whether the rule is currently active
- `applies_monthly` (BOOLEAN): Whether deduction is applied monthly
- `applies_annually` (BOOLEAN): Whether deduction is applied annually
- `is_optional` (BOOLEAN): Whether deduction is optional or mandatory
- `priority` (INTEGER): Deduction priority (higher priority applied first)
- `created_by` (UUID): User who created the rule
- `created_at` (TIMESTAMPTZ): Record creation timestamp
- `modified_by` (UUID): User who last modified the rule
- `modified_at` (TIMESTAMPTZ): Last modification timestamp

#### `deductions`

Records actual deductions applied to students' stipends.

**Columns:**

- `id` (UUID): Primary key
- `student_id` (UUID): Foreign key to students table
- `deduction_rule_id` (UUID): Foreign key to deduction_rules
- `stipend_id` (UUID): Foreign key to stipends
- `amount` (DECIMAL): Deduction amount
- `deduction_type` (VARCHAR): Type of deduction
- `description` (TEXT): Deduction description
- `deduction_date` (TIMESTAMPTZ): When deduction was applied
- `processing_status` (VARCHAR): Status (Pending, Approved, Processed, Rejected)
- `approved_by` (UUID): User who approved the deduction
- `approval_date` (TIMESTAMPTZ): When deduction was approved
- `rejection_reason` (TEXT): Reason for rejection if rejected
- `transaction_id` (UUID): Foreign key to transaction table
- `created_at` (TIMESTAMPTZ): Record creation timestamp
- `modified_at` (TIMESTAMPTZ): Last modification timestamp

## Models

### Stipend Model

```go
type Stipend struct {
    ID              uuid.UUID
    StudentID       uuid.UUID
    Amount          float64
    StipendType     string
    PaymentDate     *time.Time
    PaymentStatus   string
    PaymentMethod   string
    JournalNumber   string
    TransactionID   *uuid.UUID
    Notes           string
    CreatedAt       time.Time
    ModifiedAt      time.Time
}
```

### DeductionRule Model

```go
type DeductionRule struct {
    ID                         uuid.UUID
    RuleName                   string
    DeductionType              string
    Description                string
    BaseAmount                 float64
    MaxDeductionAmount         float64
    MinDeductionAmount         float64
    IsApplicableToFullScholar  bool
    IsApplicableToSelfFunded   bool
    IsActive                   bool
    AppliesMonthly             bool
    AppliesAnnually            bool
    IsOptional                 bool
    Priority                   int
    CreatedBy                  *uuid.UUID
    CreatedAt                  time.Time
    ModifiedBy                 *uuid.UUID
    ModifiedAt                 time.Time
}
```

### Deduction Model

```go
type Deduction struct {
    ID                 uuid.UUID
    StudentID          uuid.UUID
    DeductionRuleID    uuid.UUID
    StipendID          uuid.UUID
    Amount             float64
    DeductionType      string
    Description        string
    DeductionDate      time.Time
    ProcessingStatus   string
    ApprovedBy         *uuid.UUID
    ApprovalDate       *time.Time
    RejectionReason    string
    TransactionID      *uuid.UUID
    CreatedAt          time.Time
    ModifiedAt         time.Time
}
```

## Setup & Installation

1. **Initialize the service:**

   ```bash
   cd services/finance_service
   go mod tidy
   ```

2. **Set environment variables:**

   ```bash
   export DATABASE_URL="postgresql://user:password@localhost:5432/rub_student_portal"
   export PORT=8084
   ```

3. **Run the service:**
   ```bash
   go run main.go
   ```

## Database Migrations

The service automatically creates tables and indexes on startup:

- `stipends` table with indexes on student_id, payment_status, stipend_type, created_at
- `deduction_rules` table with indexes on rule_name, is_active, deduction_type
- `deductions` table with indexes on student_id, stipend_id, processing_status, deduction_type, deduction_rule_id

## Indexes

Created for performance optimization:

- Stipends: student_id, payment_status, stipend_type, created_at
- Deductions: student_id, stipend_id, processing_status, deduction_type, deduction_rule_id

## Testing

### gRPC Integration Tests

The Finance Service includes comprehensive gRPC integration tests:

**Documentation:**

- See [TEST_DOCUMENTATION.md](./TEST_DOCUMENTATION.md) for detailed test documentation
- See [TEST_QUICK_REFERENCE.md](./TEST_QUICK_REFERENCE.md) for quick reference

**Running Tests:**

```bash
# Start the database
docker-compose up postgres

# In another terminal, start the gRPC server
export DATABASE_URL="postgresql://postgres:postgres@localhost:5434/rub_student_portal?sslmode=disable"
./finance_service

# In another terminal, run tests
go test ./internal/grpc -v
```

## Validation Rules (Task 2.1 Implementation)

### ✅ Non-Negative Amount Validation

All financial amounts are validated to be non-negative:

**Rules:**

- Stipend amounts must be ≥ 0
- Deduction amounts must be ≥ 0
- Rule base amounts must be ≥ 0
- Warnings for amounts ≥ 100 million

**Validation Method:**

```go
validationService := NewValidationService()
result := validationService.ValidateAmount(5000.50, "Deduction amount")
if !result.IsValid {
    // Handle validation errors
}
```

### ✅ Deduction Limit Validation

Deductions are validated against rule-defined limits:

**Rules:**

- Each deduction must be ≥ rule's minimum amount
- Each deduction must be ≤ rule's maximum amount
- Deduction rules must be active to apply
- Total deductions cannot exceed stipend (by default)

**Validation Example:**

```go
result := validationService.ValidateDeductionAmount(
    3000.0,           // deduction amount
    ruleID,           // deduction rule ID
    "Deduction amount",
)
```

### ✅ Stipend Validation

Comprehensive stipend input validation:

**Validation Rules:**

- Student ID must be valid (non-nil UUID)
- Stipend type must be valid (full-scholarship, self-funded, partial)
- Amount must be non-negative
- Amount cannot exceed 10 million
- Journal number must be unique
- Journal number max length: 255 characters

**Validation Example:**

```go
result := validationService.ValidateStipendInput(
    studentID,          // UUID
    "full-scholarship", // type
    100000.0,           // amount
    "JN-STI-2024-001",  // journal number
)
```

### ✅ Deduction Rule Validation

Deduction rules are validated on creation:

**Validation Rules:**

- Rule name required, max 100 characters
- Deduction type required
- Base, min, and max amounts non-negative
- Minimum ≤ Maximum
- Base amount ≤ Maximum

**Validation Example:**

```go
result := validationService.ValidateDeductionRuleInput(
    "Hostel Fee",  // ruleName
    "hostel",      // deductionType
    5000.0,        // baseAmount
    1000.0,        // minAmount
    10000.0,       // maxAmount
)
```

### ✅ Total Deduction Validation

Total deductions are validated against stipend amount:

**Rules:**

- Total deductions cannot exceed stipend (by default)
- Warnings if deductions exceed 80% of stipend
- Can optionally allow negative net amounts

**Validation Example:**

```go
result := validationService.ValidateTotalDeductionAgainstStipend(
    4000.0,   // total deductions
    5000.0,   // stipend amount
    false,    // don't allow exceed
)
```

### Validation Results Structure

All validation methods return consistent `ValidationResult`:

```go
type ValidationResult struct {
    IsValid  bool       // true if all validations passed
    Errors   []string   // validation errors (block operation)
    Warnings []string   // warnings (informational)
}
```

**Error Handling:**

```go
if !result.IsValid {
    return fmt.Errorf(validationService.FormatValidationError(result))
}

// Log warnings
if len(result.Warnings) > 0 {
    log.Printf("%s", validationService.FormatValidationWarnings(result))
}
```

### Validation Constants

All validation limits are defined in `services/validation_constants.go`:

- `MaxStipendAmount`: 10,000,000
- `MaxDeductionAmount`: 100,000,000 (warning threshold)
- `DeductionPercentageWarning`: 80%
- `MaxRuleNameLength`: 100 characters
- `MaxJournalNumberLen`: 255 characters

For complete validation documentation, see [VALIDATION_RULES.md](./VALIDATION_RULES.md).

**Test Coverage:**

- ✅ 7 Deduction Service tests
- ✅ 6 Stipend Service tests
- ✅ 13 tests total - All passing
- ✅ ~0.35s execution time

**Key Tests:**

1. Deduction rule creation and retrieval
2. Deduction rule listing with pagination
3. Deduction creation and application
4. Stipend calculations (with deductions, monthly, annual)
5. Stipend creation and status updates
6. Error handling and edge cases

For more details on test setup, utilities, best practices, and troubleshooting, see the complete test documentation.

- Deduction Rules: rule_name, is_active, deduction_type
