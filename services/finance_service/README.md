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

1. Set environment variables:

   ```bash
   export DATABASE_URL="postgresql://user:password@host:port/rub_student_portal"
   export PORT=8084
   export GRPC_PORT=50052
   ```

2. Run the service:

   ```bash
   go run main.go
   ```

3. Service will be available at `http://localhost:8084` (REST) and `localhost:50052` (gRPC)

## Key Features (Task 1.2 + Money Transfer Implementation)

### Stipend Types

- **Full-Scholarship Students**: Limited deductions (configurable per rule)
- **Self-Funded Students**: More deductions apply by default

### ✅ Stipend Calculation

1. Base amount (annual or monthly)
2. Fetch applicable deduction rules for student type
3. Sort by priority (highest first)
4. Apply each deduction in order
5. Cap deductions to remaining stipend
6. Return base, deductions, and net amounts

### HTTP Endpoints - Stipends & Deductions

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
- `POST /api/transfers/{transactionID}/retry`: Retry a failed transfer

### Money Transfer Endpoints (NEW!)

- `POST /api/transfers/initiate`: Initiate a money transfer to student's bank
- `POST /api/transfers/{transactionID}/process`: Process pending transfer via payment gateway
- `GET /api/transfers/{transactionID}/status`: Check transfer status
- `GET /api/stipends/{stipendID}/transactions`: Get all transactions for a stipend
- `GET /api/students/{studentID}/transactions`: Get all transactions for a student
- `POST /api/transfers/{transactionID}/cancel`: Cancel a pending/processing transfer
- `POST /api/transfers/{transactionID}/retry`: Retry a failed transfer

### Search & Filter Endpoints (NEW!)

- `GET /api/stipends/search`: Search stipends with filters (student_id, status, type, amount range, date range)
- `GET /api/deduction-rules/search`: Search deduction rules (name, type, is_active)
- `GET /api/transactions/search`: Search transactions (student_id, stipend_id, status, amount range, date range)

### Report Generation Endpoints (NEW!)

- `GET /api/reports/disbursement`: Generate disbursement summary report
- `GET /api/reports/deductions`: Generate deduction summary report
- `GET /api/reports/transactions`: Generate transaction summary report
- `GET /api/reports/export/stipends`: Export stipends to CSV
- `GET /api/reports/export/deductions`: Export deductions to CSV
- `GET /api/reports/export/transactions`: Export transactions to CSV

### Audit Log Endpoints (NEW!)

- `GET /api/audit-logs`: Get all audit logs with optional filters (action, entity_type, finance_officer, status, date range)
- `GET /api/audit-logs/entity/{entityType}/{entityID}`: Get audit logs for specific entity
- `GET /api/audit-logs/officer/{officerID}`: Get audit logs for specific finance officer

### Deduction Management

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

#### `transactions`

Records money transfer transactions for stipend disbursement.

**Columns:**

- `id` (UUID): Primary key
- `stipend_id` (UUID): Foreign key to stipends
- `student_id` (UUID): Foreign key to students
- `amount` (DECIMAL): Transfer amount
- `source_account` (VARCHAR): College/Institution account
- `destination_account` (VARCHAR): Student's bank account
- `destination_bank` (VARCHAR): Student's bank name
- `transaction_type` (VARCHAR): Type (STIPEND, REFUND, etc)
- `status` (VARCHAR): Status (PENDING, PROCESSING, SUCCESS, FAILED, CANCELLED)
- `payment_method` (VARCHAR): Payment method (BANK_TRANSFER, E_PAYMENT)
- `reference_number` (VARCHAR): Unique reference from payment gateway
- `error_message` (TEXT): Error details if failed
- `remarks` (TEXT): Additional remarks
- `initiated_at` (TIMESTAMPTZ): When transfer was initiated
- `processed_at` (TIMESTAMPTZ): When transfer was processed
- `completed_at` (TIMESTAMPTZ): When transfer was completed
- `created_at` (TIMESTAMPTZ): Record creation timestamp
- `modified_at` (TIMESTAMPTZ): Last modification timestamp

#### `audit_logs`

Tracks all financial operations for compliance and transparency.

**Columns:**

- `id` (UUID): Primary key
- `action` (VARCHAR): Action type (CREATE, UPDATE, DELETE, VIEW)
- `entity_type` (VARCHAR): Entity type (STIPEND, DEDUCTION_RULE, TRANSACTION)
- `entity_id` (UUID): ID of the entity being acted upon
- `finance_officer` (VARCHAR): Email/ID of the finance officer
- `description` (TEXT): Action description
- `old_values` (JSONB): Previous data for updates
- `new_values` (JSONB): New data
- `status` (VARCHAR): Status (SUCCESS, FAILED)
- `error_message` (TEXT): Error details if failed
- `ip_address` (VARCHAR): IP address of requester
- `user_agent` (TEXT): User agent string
- `created_at` (TIMESTAMPTZ): Record creation timestamp
- `updated_at` (TIMESTAMPTZ): Last update timestamp

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

### Transaction Model

```go
type Transaction struct {
    ID                 uuid.UUID
    StipendID          uuid.UUID
    StudentID          uuid.UUID
    Amount             float64
    SourceAccount      string
    DestinationAccount string
    DestinationBank    string
    TransactionType    string
    Status             string
    PaymentMethod      string
    ReferenceNumber    sql.NullString
    ErrorMessage       string
    Remarks            string
    InitiatedAt        time.Time
    ProcessedAt        *time.Time
    CompletedAt        *time.Time
    CreatedAt          time.Time
    ModifiedAt         time.Time
}
```

### AuditLog Model

```go
type AuditLog struct {
    ID            string
    Action        string
    EntityType    string
    EntityID      string
    FinanceOfficer string
    Description   string
    OldValues     string // JSONB
    NewValues     string // JSONB
    Status        string
    ErrorMessage  string
    IPAddress     string
    UserAgent     string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

## Services

### Business Logic Services

1. **StipendService** - Stipend creation and management
2. **DeductionService** - Apply deductions to stipends
3. **DeductionRuleService** - Manage deduction rules
4. **TransferService** - Handle money transfers and transactions
5. **ReportService** - Generate financial reports and export CSV
6. **SearchService** - Search and filter stipends, deductions, transactions
7. **AuditService** - Track and log all financial operations
8. **ValidationService** - Validate financial inputs and amounts
9. **ErrorLogger** - Centralized error logging
10. **BankingClient** - gRPC client for banking service integration
11. **StudentClient** - gRPC client for student service integration

### gRPC Services

- **StipendService** (8 RPC methods) - Calculate, validate, and manage stipends
- **DeductionService** (8 RPC methods) - Create and apply deduction rules

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
- `transactions` table with indexes on stipend_id, student_id, status, created_at
- `audit_logs` table with indexes on action, entity_type, finance_officer, created_at

## Indexes

Created for performance optimization:

- Stipends: student_id, payment_status, stipend_type, created_at
- Deductions: student_id, stipend_id, processing_status, deduction_type, deduction_rule_id
- Transactions: stipend_id, student_id, status, created_at
- AuditLogs: action, entity_type, finance_officer, created_at

## Testing

### Comprehensive Test Coverage

**Services tested:**

- 11 Business logic services (Stipend, Deduction, Transfer, Report, Search, Audit, Validation, etc.)
- 2 gRPC services (Stipend service, Deduction service) with 16 RPC methods total
- 7 HTTP handlers with 30 REST endpoints

**Unit Tests:**

- Stipend service tests
- Deduction service tests
- Transfer service tests (8 tests for transaction lifecycle)
- Validation service tests

### gRPC Integration Tests

The Finance Service includes comprehensive gRPC integration tests.

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

### Transaction Lifecycle Tests

Comprehensive tests for money transfer workflows:

**Test Cases:**

1. `TestTransactionCreation` - Basic transaction creation
2. `TestTransactionStatusUpdate` - Status transitions (Pending → Processing)
3. `TestTransactionSuccess` - Successful transaction completion
4. `TestTransactionFailed` - Failed transaction handling
5. `TestTransactionQuery` - Query by stipend and student
6. `TestTransactionCancellation` - Canceling pending/processing transactions

## Validation Rules (Task 2.1 Implementation)

### Non-Negative Amount Validation

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

### Deduction Limit Validation

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

### Stipend Validation

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

### Deduction Rule Validation

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

### Total Deduction Validation

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

**Test Coverage:**

- 7 Deduction Service tests
- 6 Stipend Service tests
- 13 tests total - All passing
- ~0.35s execution time

**Key Tests:**

1. Deduction rule creation and retrieval
2. Deduction rule listing with pagination
3. Deduction creation and application
4. Stipend calculations (with deductions, monthly, annual)
5. Stipend creation and status updates
6. Error handling and edge cases

### Handler Layer (7 handlers, 30 REST endpoints)

- **StipendHandler** - Stipend CRUD operations (8 endpoints)
- **DeductionHandler** - Deduction rule management (3 endpoints)
- **TransferHandler** - Money transfer operations (5 endpoints)
- **SearchHandler** - Search and filtering (3 endpoints)
- **ReportHandler** - Report generation and CSV export (6 endpoints)
- **AuditHandler** - Audit log queries (3 endpoints)
- **Health check** - Service health (2 endpoints)

### Service Layer (11 services)

Business logic implementation with database integration and validation:

- Stipend calculations and management
- Deduction rule application
- Money transfer processing
- Report generation
- Advanced searching and filtering
- Audit logging for compliance
- Input validation and error handling
- gRPC client integration (Banking, Student services)

### gRPC Services (2 services, 16 RPC methods)

Inter-service communication:

- **StipendService**: Calculate, validate, and manage stipends via gRPC
- **DeductionService**: Create and apply deduction rules via gRPC

### Database Layer

- PostgreSQL with GORM ORM
- 5 tables (stipends, deduction_rules, deductions, transactions, audit_logs)
- Automatic migrations and seeding
- Performance-optimized indexes

## Features Summary

### Core Features

- ✅ Stipend calculation with deductions
- ✅ Configurable deduction rules by student type
- ✅ Priority-based deduction ordering
- ✅ Payment tracking and status management

### Money Transfer Module

- ✅ Initiate transfers to student bank accounts
- ✅ Payment gateway integration (simulated)
- ✅ Transfer status tracking
- ✅ Retry failed transfers
- ✅ Cancel pending transfers
- ✅ Transaction audit trail

### Search & Filtering

- ✅ Search stipends with multi-filter support
- ✅ Search deduction rules
- ✅ Search transactions with advanced filters
- ✅ Pagination support

### Reporting

- ✅ Disbursement summary reports
- ✅ Deduction analysis reports
- ✅ Transaction summary reports
- ✅ CSV export for stipends, deductions, transactions

### Audit & Compliance

- ✅ Complete audit logging of all operations
- ✅ Track actions by finance officer
- ✅ Store old and new values for updates
- ✅ IP address and user agent tracking
- ✅ Filter audit logs by entity, action, date range

### Validation

- ✅ Amount validation (non-negative, reasonable bounds)
- ✅ Deduction limit validation
- ✅ Stipend input validation
- ✅ Deduction rule validation
- ✅ Total deduction validation against stipend

## Port Configuration

- **REST API**: 8084
- **gRPC Server**: 50052
- **Database**: PostgreSQL (Render.com cloud)

## Cloud Database

The Finance Service uses a **PostgreSQL database hosted on Render.com** for reliability and scalability:

**Database Infrastructure:**

- **Provider**: Render.com PostgreSQL instance
- **Connection**: Secure SSL/TLS encrypted connections
- **Backups**: Automatic daily backups
- **Availability**: 99.9% uptime SLA
- **Scalability**: Automatic resource scaling
- **Replication**: Multi-region backup replication
- **Monitoring**: Real-time database metrics and alerts

**Benefits:**

- No local database setup required
- Production-ready configuration
- Automatic maintenance and patches
- Data persistence and recovery
- Geographic redundancy

Connection is managed via the `DATABASE_URL` environment variable.

## Environment Variables

```bash
DATABASE_URL=postgresql://user:password@host:port/rub_student_portal
PORT=8084
GRPC_PORT=50052
ENVIRONMENT=production
```
