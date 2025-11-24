# Finance Service - Stipend Calculation & Deduction System

## Overview

The Finance Service manages stipend calculations and deductions for students. It supports:

- Full-scholarship and self-funded student stipend management
- Configurable deduction rules (hostel, electricity, mess fees, etc.)
- Deduction tracking and audit logging
- Payment processing and status management

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
- Deduction Rules: rule_name, is_active, deduction_type
