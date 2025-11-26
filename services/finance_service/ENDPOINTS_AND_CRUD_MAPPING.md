# Finance Service - Endpoints & CRUD Functions Mapping

## 📍 Project Structure

```
finance_service/
├── handlers/              # HTTP Request Handlers
│   ├── stipend_handler.go
│   ├── deduction_handler.go
│   ├── transfer_handler.go
│   └── deduction_handler_test.go
├── services/              # Business Logic & CRUD Operations
│   ├── stipend_service.go
│   ├── deduction_rule_service.go
│   ├── transfer_service.go
│   ├── deduction_service.go
│   └── *_test.go files
├── models/                # Database Models
│   ├── stipend.go
│   ├── deduction.go
│   ├── deduction_rule.go
│   └── transaction.go
├── database/              # Database Connection & Migrations
│   └── db.go
└── main.go               # Router & Server Setup
```

---

## 🔄 CRUD Operations Overview

### 1. STIPEND SERVICE (`stipend_service.go`)

**CREATE Operations:**

- `CreateStipendForStudent()` - Create new stipend for a student

**READ Operations:**

- `GetStipendByID()` - Get single stipend by ID
- `GetStudentStipends()` - Get all stipends for a student (with pagination)
- `GetStudentStipendsWithPagination()` - Enhanced pagination support
- `GetStipendDeductions()` - Get all deductions for a stipend

**UPDATE Operations:**

- `UpdateStipendPaymentStatus()` - Update payment status (Pending/Processed/Failed)
- `UpdateStipendPaymentStatusWithReturn()` - Update with return value

**DELETE Operations:**

- None directly (handled via cascading deletes in database)

---

### 2. DEDUCTION RULE SERVICE (`deduction_rule_service.go`)

**CREATE Operations:**

- `CreateRule()` - Create new deduction rule
- `CreateDeductionRule()` - Alternative create method

**READ Operations:**

- `GetRuleByID()` - Get single rule by ID
- `ListActiveRules()` - List only active rules with pagination
- `ListAllRules()` - List all rules (active & inactive) with pagination
- `ListRulesByType()` - Filter rules by deduction type with pagination
- `GetApplicableRules()` - Get rules applicable to scholar type

**UPDATE Operations:**

- `UpdateRule()` - Update rule details
- `UpdateDeductionRule()` - Alternative update method

**DELETE Operations:**

- `DeleteRule()` - Delete a deduction rule

---

### 3. TRANSFER SERVICE (`transfer_service.go`)

**CREATE Operations:**

- `InitiateTransfer()` - Create new money transfer transaction

**READ Operations:**

- `GetTransactionStatus()` - Get status of a transaction
- `GetTransactionsByStipend()` - Get all transactions for a stipend
- `GetTransactionsByStudent()` - Get all transactions for a student

**UPDATE Operations:**

- None directly implemented (status updates handled internally)

**DELETE Operations:**

- None

---

## 📡 API Endpoints Mapping

### Stipend Endpoints

| HTTP Method | Endpoint                                   | Handler Function                   | CRUD Type | Service Function                           |
| ----------- | ------------------------------------------ | ---------------------------------- | --------- | ------------------------------------------ |
| POST        | `/api/stipends`                            | `CreateStipend()`                  | CREATE    | `CreateStipendForStudent()`                |
| GET         | `/api/stipends/{stipendID}`                | `GetStipend()`                     | READ      | `GetStipendByID()`                         |
| GET         | `/api/students/{studentID}/stipends`       | `GetStudentStipends()`             | READ      | `GetStudentStipends()`                     |
| PATCH       | `/api/stipends/{stipendID}/payment-status` | `UpdateStipendPaymentStatus()`     | UPDATE    | `UpdateStipendPaymentStatus()`             |
| GET         | `/api/stipends/{stipendID}/deductions`     | `GetStipendDeductions()`           | READ      | `GetStipendDeductions()`                   |
| POST        | `/api/stipends/calculate`                  | `CalculateStipendWithDeductions()` | READ/CALC | `CreateStipendForStudent()` + calculations |
| POST        | `/api/stipends/calculate/monthly`          | `CalculateMonthlyStipend()`        | READ/CALC | Custom calculation                         |
| POST        | `/api/stipends/calculate/annual`           | `CalculateAnnualStipend()`         | READ/CALC | Custom calculation                         |

---

### Deduction Rule Endpoints

| HTTP Method | Endpoint                        | Handler Function        | CRUD Type | Service Function |
| ----------- | ------------------------------- | ----------------------- | --------- | ---------------- |
| POST        | `/api/deduction-rules`          | `CreateDeductionRule()` | CREATE    | `CreateRule()`   |
| GET         | `/api/deduction-rules`          | `ListDeductionRules()`  | READ      | `ListAllRules()` |
| GET         | `/api/deduction-rules/{ruleID}` | `GetDeductionRule()`    | READ      | `GetRuleByID()`  |
| PUT         | `/api/deduction-rules/{ruleID}` | `UpdateDeductionRule()` | UPDATE    | `UpdateRule()`   |
| DELETE      | `/api/deduction-rules/{ruleID}` | `DeleteDeductionRule()` | DELETE    | `DeleteRule()`   |

---

### Transfer/Transaction Endpoints

| HTTP Method | Endpoint                                 | Handler Function             | CRUD Type | Service Function             |
| ----------- | ---------------------------------------- | ---------------------------- | --------- | ---------------------------- |
| POST        | `/api/transfers/initiate`                | `InitiateTransfer()`         | CREATE    | `InitiateTransfer()`         |
| GET         | `/api/stipends/{stipendID}/transactions` | `GetTransactionsByStipend()` | READ      | `GetTransactionsByStipend()` |
| GET         | `/api/students/{studentID}/transactions` | `GetTransactionsByStudent()` | READ      | `GetTransactionsByStudent()` |

---

## 📊 Data Flow Example: Creating a Stipend

```
1. HTTP Request (POST /api/v1/finance/stipends)
   ↓
2. Handler: StipendHandler.CreateStipend()
   ↓
3. Parse JSON Request
   ↓
4. Validation
   ↓
5. Service: StipendService.CreateStipendForStudent()
   ↓
6. Database: GORM creates record in 'stipends' table
   ↓
7. Response with created stipend data
```

---

## 🗂️ Service Layer Organization

### `stipend_service.go` - **262+ lines**

Handles all stipend-related operations and deduction rule management

### `deduction_rule_service.go` - **450+ lines**

Dedicated service for deduction rule CRUD with validation and error logging

### `transfer_service.go` - **180+ lines**

Handles money transfer transactions and payment processing

### `deduction_service.go` - **150+ lines**

Handles deduction calculations and application

### Supporting Services:

- `validation_service.go` - Input validation and business rules
- `error_logger.go` - Error tracking and logging
- `banking_client.go` - Banking service integration
- `student_client.go` - Student data integration
- `user_client.go` - User data integration

---

## ✅ Working Status

| Component                | Status       | Details                                     |
| ------------------------ | ------------ | ------------------------------------------- |
| **Stipend CRUD**         | ✅ Working   | Create, Read, Update operations implemented |
| **Deduction Rules CRUD** | ✅ Working   | Full CRUD with validation                   |
| **Transfers**            | ✅ Working   | Transaction initiation and status tracking  |
| **Database**             | ✅ Connected | Cloud database (Render)                     |
| **Handlers**             | ✅ Working   | All HTTP handlers implemented               |
| **Services**             | ✅ Working   | All business logic implemented              |
| **Models**               | ✅ Working   | All database models defined                 |

---

## 🚀 How to Test Each CRUD Operation

### Create Stipend

```bash
curl -X POST http://localhost:8084/api/stipends \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": 25000.00,
    "stipend_type": "full-scholarship",
    "payment_method": "Bank_transfer",
    "journal_number": "JNL-2025-001",
    "notes": "Full scholarship for academic year 2025"
  }'
```

### Read Stipend

```bash
curl -X GET http://localhost:8084/api/stipends/550e8400-e29b-41d4-a716-446655440001 \
  -H "Content-Type: application/json"
```

### Update Stipend Status

```bash
curl -X PATCH http://localhost:8084/api/stipends/550e8400-e29b-41d4-a716-446655440001/payment-status \
  -H "Content-Type: application/json" \
  -d '{
    "payment_status": "Processed"
  }'
```

### Get Student Stipends

```bash
curl -X GET http://localhost:8084/api/students/550e8400-e29b-41d4-a716-446655440000/stipends \
  -H "Content-Type: application/json"
```

### Get Stipend Deductions

```bash
curl -X GET http://localhost:8084/api/stipends/550e8400-e29b-41d4-a716-446655440001/deductions \
  -H "Content-Type: application/json"
```

### Create Deduction Rule

```bash
curl -X POST http://localhost:8084/api/deduction-rules \
  -H "Content-Type: application/json" \
  -d '{
    "rule_name": "Hostel Fee 2025",
    "deduction_type": "hostel",
    "description": "Monthly hostel accommodation fee",
    "base_amount": 5000.00,
    "max_deduction_amount": 5000.00,
    "min_deduction_amount": 1000.00,
    "is_applicable_to_full_scholar": false,
    "is_applicable_to_self_funded": true,
    "is_active": true,
    "applies_monthly": true,
    "applies_annually": false,
    "is_optional": false,
    "priority": 1
  }'
```

### List All Deduction Rules

```bash
curl -X GET http://localhost:8084/api/deduction-rules \
  -H "Content-Type: application/json"
```

### Get Specific Deduction Rule

```bash
curl -X GET http://localhost:8084/api/deduction-rules/550e8400-e29b-41d4-a716-446655440002 \
  -H "Content-Type: application/json"
```

### Update Deduction Rule

```bash
curl -X PUT http://localhost:8084/api/deduction-rules/550e8400-e29b-41d4-a716-446655440002 \
  -H "Content-Type: application/json" \
  -d '{
    "base_amount": 6000.00,
    "max_deduction_amount": 6000.00,
    "is_active": true
  }'
```

### Delete Deduction Rule

```bash
curl -X DELETE http://localhost:8084/api/deduction-rules/550e8400-e29b-41d4-a716-446655440002 \
  -H "Content-Type: application/json"
```

### Calculate Stipend with Deductions

```bash
curl -X POST http://localhost:8084/api/stipends/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "stipend_id": "550e8400-e29b-41d4-a716-446655440001",
    "deduction_rule_ids": [
      "550e8400-e29b-41d4-a716-446655440002",
      "550e8400-e29b-41d4-a716-446655440003"
    ]
  }'
```

### Calculate Monthly Stipend

```bash
curl -X POST http://localhost:8084/api/stipends/calculate/monthly \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "base_amount": 25000.00
  }'
```

### Calculate Annual Stipend

```bash
curl -X POST http://localhost:8084/api/stipends/calculate/annual \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "base_amount": 300000.00
  }'
```

### Initiate Transfer

```bash
curl -X POST http://localhost:8084/api/transfers/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "stipend_id": "550e8400-e29b-41d4-a716-446655440001",
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "amount": 20000.00,
    "destination_account": "1234567890",
    "destination_bank": "BDB",
    "payment_method": "BANK_TRANSFER"
  }'
```

### Get Transactions by Stipend

```bash
curl -X GET http://localhost:8084/api/stipends/550e8400-e29b-41d4-a716-446655440001/transactions \
  -H "Content-Type: application/json"
```

### Get Transactions by Student

```bash
curl -X GET http://localhost:8084/api/students/550e8400-e29b-41d4-a716-446655440000/transactions \
  -H "Content-Type: application/json"
```

### Health Check

```bash
curl -X GET http://localhost:8084/health \
  -H "Content-Type: application/json"
```

---

## 📝 Notes

- All CRUD operations are **fully implemented** and working
- Service layer handles all business logic
- Handlers manage HTTP communication
- Database operations use GORM ORM
- Error handling and logging included throughout
- All services are **tested** with unit tests

---

## 🔍 SEARCH & FILTER ENDPOINTS

### Search Stipends

```bash
GET /api/search/stipends
```

**Query Parameters:**

- `student_id` - Filter by student ID
- `status` - Filter by payment status (Pending, Processed, Failed)
- `stipend_type` - Filter by type (full-scholarship, self-funded)
- `payment_status` - Filter by payment status
- `start_date` - RFC3339 format (e.g., 2025-01-01T00:00:00Z)
- `end_date` - RFC3339 format
- `min_amount` - Minimum stipend amount
- `max_amount` - Maximum stipend amount
- `limit` - Results per page (default: 10, max: 100)
- `offset` - Page offset (default: 0)

**Example:**

```bash
curl -s "http://localhost:8084/api/search/stipends?student_id=550e8400-e29b-41d4-a716-446655440000&limit=10&offset=0" | jq .
```

### Search Deduction Rules

```bash
GET /api/search/deduction-rules
```

**Query Parameters:**

- `rule_name` - Search rule name (case-insensitive partial match)
- `deduction_type` - Filter by deduction type
- `is_active` - Filter by active status (true/false)
- `limit` - Results per page (default: 10, max: 100)
- `offset` - Page offset (default: 0)

**Example:**

```bash
curl -s "http://localhost:8084/api/search/deduction-rules?deduction_type=hostel&is_active=true" | jq .
```

### Search Transactions

```bash
GET /api/search/transactions
```

**Query Parameters:**

- `student_id` - Filter by student ID
- `stipend_id` - Filter by stipend ID
- `status` - Filter by transaction status (PENDING, SUCCESS, FAILED)
- `transaction_type` - Filter by transaction type
- `start_date` - RFC3339 format
- `end_date` - RFC3339 format
- `min_amount` - Minimum transaction amount
- `max_amount` - Maximum transaction amount
- `limit` - Results per page (default: 10, max: 100)
- `offset` - Page offset (default: 0)

**Example:**

```bash
curl -s "http://localhost:8084/api/search/transactions?student_id=550e8400-e29b-41d4-a716-446655440000&status=SUCCESS" | jq .
```

---

## 📊 REPORT GENERATION ENDPOINTS

### Disbursement Report

```bash
GET /api/reports/disbursement
```

**Query Parameters:**

- `start_date` - RFC3339 format (optional)
- `end_date` - RFC3339 format (optional)

**Response Example:**

```json
{
  "total_stipends": 1,
  "total_amount": 25000,
  "pending_count": 1,
  "processed_count": 0,
  "failed_count": 0,
  "average_amount": 25000,
  "min_amount": 25000,
  "max_amount": 25000,
  "generated_at": "2025-11-26T21:53:51.633591723+06:00",
  "report_period": ""
}
```

**Use Case:** Finance Officer gets overview of all disbursements

### Deduction Report

```bash
GET /api/reports/deductions
```

**Response:** Array of deduction types with totals and statistics

**Use Case:** Finance Officer analyzes deduction patterns

### Transaction Report

```bash
GET /api/reports/transactions
```

**Query Parameters:**

- `start_date` - RFC3339 format (optional)
- `end_date` - RFC3339 format (optional)

**Response Example:**

```json
{
  "total_transactions": 0,
  "successful_count": 0,
  "pending_count": 0,
  "failed_count": 0,
  "total_amount": 0,
  "average_amount": 0,
  "generated_at": "2025-11-26T21:53:51.633591723+06:00",
  "report_period": ""
}
```

**Use Case:** Finance Officer monitors transaction success rates

### Export Stipends to CSV

```bash
GET /api/reports/export/stipends
```

**Query Parameters:**

- `start_date` - RFC3339 format (optional)
- `end_date` - RFC3339 format (optional)

**Response:** CSV file download with headers: ID, Student ID, Amount, Stipend Type, Payment Status, Journal Number, Created At, Modified At

### Export Deductions to CSV

```bash
GET /api/reports/export/deductions
```

**Query Parameters:**

- `start_date` - RFC3339 format (optional)
- `end_date` - RFC3339 format (optional)

**Response:** CSV file with deduction details

### Export Transactions to CSV

```bash
GET /api/reports/export/transactions
```

**Query Parameters:**

- `start_date` - RFC3339 format (optional)
- `end_date` - RFC3339 format (optional)

**Response:** CSV file with transaction details

---

## 📋 AUDIT LOGGING ENDPOINTS

### Get All Audit Logs

```bash
GET /api/audit-logs
```

**Query Parameters:**

- `action` - Filter by action (CREATE, UPDATE, DELETE)
- `entity_type` - Filter by entity type (STIPEND, DEDUCTION_RULE, TRANSACTION)
- `finance_officer` - Filter by officer who made the action
- `status` - Filter by log status (SUCCESS, FAILED)
- `start_date` - RFC3339 format
- `end_date` - RFC3339 format
- `limit` - Results per page (default: 10, max: 100)
- `offset` - Page offset (default: 0)

**Example:**

```bash
curl -s "http://localhost:8084/api/audit-logs?action=CREATE&entity_type=STIPEND" | jq .
```

### Get Audit Logs by Entity

```bash
GET /api/audit-logs/{entity_type}/{entity_id}
```

**Path Parameters:**

- `entity_type` - Entity type (STIPEND, DEDUCTION_RULE, TRANSACTION)
- `entity_id` - Entity UUID

**Response:** Complete audit history for the specific entity

**Use Case:** Trace all changes made to a specific stipend or deduction rule

### Get Audit Logs by Finance Officer

```bash
GET /api/audit-logs/officer/{officer}
```

**Path Parameters:**

- `officer` - Finance officer email/ID

**Query Parameters:**

- `limit` - Results per page (default: 10)
- `offset` - Page offset (default: 0)

**Response:** All actions performed by this officer

**Use Case:** Monitor specific officer's activities for compliance

---

## 📝 New Features Summary

- **NEW:** Search, Filter, and Pagination support all major entities
- **NEW:** Report generation available in JSON and CSV formats
- **NEW:** Comprehensive audit logging for compliance and transparency
- Audit logs capture who did what, when, and whether it succeeded
