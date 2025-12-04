# Finance Service - API Reference

## Base URL

```
http://localhost:8084/api
```

## Authentication

Currently, no authentication is required. Add authentication middleware as needed.

---

## Stipend Endpoints

### 1. Create Stipend

Creates a new stipend record for a student.

**Endpoint**: `POST /stipends`

**Request Body**:

```json
{
  "student_id": "string (UUID)",
  "stipend_type": "string (full-scholarship|self-funded)",
  "amount": "number",
  "payment_method": "string",
  "journal_number": "string",
  "notes": "string (optional)"
}
```

**Response** (201 Created):

```json
{
  "id": "string (UUID)",
  "student_id": "string (UUID)",
  "amount": "number",
  "stipend_type": "string",
  "payment_date": "string (ISO8601) | null",
  "payment_status": "string (Pending|Processed|Failed)",
  "payment_method": "string",
  "journal_number": "string",
  "notes": "string",
  "created_at": "string (ISO8601)",
  "modified_at": "string (ISO8601)"
}
```

**Example**:

```bash
curl -X POST http://localhost:8084/api/stipends \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "stipend_type": "full-scholarship",
    "amount": 50000.00,
    "payment_method": "Bank_transfer",
    "journal_number": "JN-001-2024",
    "notes": "January 2024 stipend"
  }'
```

---

### 2. Calculate Stipend with Deductions

Calculates stipend amount after applying applicable deductions without creating a record.

**Endpoint**: `POST /stipends/calculate`

**Request Body**:

```json
{
  "student_id": "string (UUID)",
  "stipend_type": "string (full-scholarship|self-funded)",
  "amount": "number",
  "payment_method": "string",
  "journal_number": "string",
  "notes": "string (optional)"
}
```

**Response** (200 OK):

```json
{
  "base_stipend_amount": "number",
  "total_deductions": "number",
  "net_stipend_amount": "number",
  "deductions": [
    {
      "rule_id": "string (UUID)",
      "rule_name": "string",
      "deduction_type": "string",
      "amount": "number",
      "description": "string",
      "is_optional": "boolean"
    }
  ]
}
```

---

### 3. Calculate Monthly Stipend

Calculates monthly stipend (annual ÷ 12) with deductions.

**Endpoint**: `POST /stipends/calculate/monthly`

**Request Body**:

```json
{
  "student_id": "string (UUID)",
  "stipend_type": "string (full-scholarship|self-funded)",
  "annual_amount": "number"
}
```

**Response** (200 OK):

```json
{
  "base_stipend_amount": "number",
  "total_deductions": "number",
  "net_stipend_amount": "number",
  "deductions": [...]
}
```

---

### 4. Calculate Annual Stipend

Calculates annual stipend with applicable deductions.

**Endpoint**: `POST /stipends/calculate/annual`

**Request Body**:

```json
{
  "student_id": "string (UUID)",
  "stipend_type": "string (full-scholarship|self-funded)",
  "annual_amount": "number"
}
```

**Response** (200 OK):

```json
{
  "base_stipend_amount": "number",
  "total_deductions": "number",
  "net_stipend_amount": "number",
  "deductions": [...]
}
```

---

### 5. Get Stipend by ID

Retrieves a specific stipend record.

**Endpoint**: `GET /stipends/{stipendID}`

**Path Parameters**:

- `stipendID` (UUID): The ID of the stipend

**Response** (200 OK):

```json
{
  "id": "string (UUID)",
  "student_id": "string (UUID)",
  "amount": "number",
  "stipend_type": "string",
  "payment_date": "string (ISO8601) | null",
  "payment_status": "string",
  "payment_method": "string",
  "journal_number": "string",
  "notes": "string",
  "created_at": "string (ISO8601)",
  "modified_at": "string (ISO8601)"
}
```

**Error** (404 Not Found):

```json
{
  "error": "Stipend not found"
}
```

---

### 6. Get Student's Stipends

Retrieves all stipends for a specific student with pagination.

**Endpoint**: `GET /students/{studentID}/stipends`

**Path Parameters**:

- `studentID` (UUID): The ID of the student

**Query Parameters**:

- `limit` (optional, default: 10): Number of records per page
- `offset` (optional, default: 0): Number of records to skip

**Response** (200 OK):

```json
{
  "stipends": [
    {
      "id": "string (UUID)",
      "student_id": "string (UUID)",
      "amount": "number",
      "stipend_type": "string",
      "payment_date": "string (ISO8601) | null",
      "payment_status": "string",
      "payment_method": "string",
      "journal_number": "string",
      "notes": "string",
      "created_at": "string (ISO8601)",
      "modified_at": "string (ISO8601)"
    }
  ],
  "total": "number",
  "limit": "number",
  "offset": "number"
}
```

**Example**:

```bash
curl -X GET "http://localhost:8084/api/students/550e8400-e29b-41d4-a716-446655440000/stipends?limit=10&offset=0"
```

---

### 7. Update Stipend Payment Status

Updates the payment status of a stipend.

**Endpoint**: `PATCH /stipends/{stipendID}/payment-status`

**Path Parameters**:

- `stipendID` (UUID): The ID of the stipend

**Request Body**:

```json
{
  "status": "string (Pending|Processed|Failed)",
  "payment_date": "string (ISO8601, optional)"
}
```

**Response** (200 OK):

```json
{
  "id": "string (UUID)",
  "student_id": "string (UUID)",
  "amount": "number",
  "stipend_type": "string",
  "payment_date": "string (ISO8601) | null",
  "payment_status": "string",
  "payment_method": "string",
  "journal_number": "string",
  "notes": "string",
  "created_at": "string (ISO8601)",
  "modified_at": "string (ISO8601)"
}
```

**Example**:

```bash
curl -X PATCH http://localhost:8084/api/stipends/f47ac10b-58cc-4372-a567-0e02b2c3d479/payment-status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Processed",
    "payment_date": "2024-01-15T10:30:00Z"
  }'
```

---

### 8. Get Stipend Deductions

Retrieves all deductions applied to a specific stipend.

**Endpoint**: `GET /stipends/{stipendID}/deductions`

**Path Parameters**:

- `stipendID` (UUID): The ID of the stipend

**Response** (200 OK):

```json
{
  "deductions": [
    {
      "id": "string (UUID)",
      "student_id": "string (UUID)",
      "deduction_rule_id": "string (UUID)",
      "stipend_id": "string (UUID)",
      "amount": "number",
      "deduction_type": "string",
      "description": "string",
      "deduction_date": "string (ISO8601)",
      "processing_status": "string",
      "approved_by": "string (UUID) | null",
      "approval_date": "string (ISO8601) | null",
      "rejection_reason": "string",
      "created_at": "string (ISO8601)",
      "modified_at": "string (ISO8601)"
    }
  ],
  "count": "number"
}
```

---

## Deduction Rule Endpoints

### 1. Create Deduction Rule

Creates a new deduction rule that can be applied to stipends.

**Endpoint**: `POST /deduction-rules`

**Request Body**:

```json
{
  "rule_name": "string",
  "deduction_type": "string",
  "description": "string",
  "base_amount": "number",
  "max_deduction_amount": "number",
  "min_deduction_amount": "number (optional, default: 0)",
  "is_applicable_to_full_scholar": "boolean",
  "is_applicable_to_self_funded": "boolean",
  "applies_monthly": "boolean",
  "applies_annually": "boolean",
  "is_optional": "boolean",
  "priority": "number"
}
```

**Response** (201 Created):

```json
{
  "id": "string (UUID)",
  "rule_name": "string",
  "deduction_type": "string",
  "description": "string",
  "base_amount": "number",
  "max_deduction_amount": "number",
  "min_deduction_amount": "number",
  "is_applicable_to_full_scholar": "boolean",
  "is_applicable_to_self_funded": "boolean",
  "is_active": "boolean",
  "applies_monthly": "boolean",
  "applies_annually": "boolean",
  "is_optional": "boolean",
  "priority": "number",
  "created_at": "string (ISO8601)",
  "modified_at": "string (ISO8601)"
}
```

**Example**:

```bash
curl -X POST http://localhost:8084/api/deduction-rules \
  -H "Content-Type: application/json" \
  -d '{
    "rule_name": "Hostel Fee",
    "deduction_type": "hostel",
    "description": "Monthly hostel charges",
    "base_amount": 3000.00,
    "max_deduction_amount": 3500.00,
    "min_deduction_amount": 2500.00,
    "is_applicable_to_full_scholar": true,
    "is_applicable_to_self_funded": true,
    "applies_monthly": true,
    "applies_annually": false,
    "is_optional": false,
    "priority": 100
  }'
```

---

### 2. Get Deduction Rule by ID

Retrieves a specific deduction rule.

**Endpoint**: `GET /deduction-rules/{ruleID}`

**Path Parameters**:

- `ruleID` (UUID): The ID of the deduction rule

**Response** (200 OK):

```json
{
  "id": "string (UUID)",
  "rule_name": "string",
  "deduction_type": "string",
  "description": "string",
  "base_amount": "number",
  "max_deduction_amount": "number",
  "min_deduction_amount": "number",
  "is_applicable_to_full_scholar": "boolean",
  "is_applicable_to_self_funded": "boolean",
  "is_active": "boolean",
  "applies_monthly": "boolean",
  "applies_annually": "boolean",
  "is_optional": "boolean",
  "priority": "number",
  "created_at": "string (ISO8601)",
  "modified_at": "string (ISO8601)"
}
```

---

### 3. List Deduction Rules

Retrieves all active deduction rules with pagination.

**Endpoint**: `GET /deduction-rules`

**Query Parameters**:

- `limit` (optional, default: 20): Number of records per page
- `offset` (optional, default: 0): Number of records to skip

**Response** (200 OK):

```json
{
  "rules": [
    {
      "id": "string (UUID)",
      "rule_name": "string",
      "deduction_type": "string",
      "description": "string",
      "base_amount": "number",
      "max_deduction_amount": "number",
      "min_deduction_amount": "number",
      "is_applicable_to_full_scholar": "boolean",
      "is_applicable_to_self_funded": "boolean",
      "is_active": "boolean",
      "applies_monthly": "boolean",
      "applies_annually": "boolean",
      "is_optional": "boolean",
      "priority": "number",
      "created_at": "string (ISO8601)",
      "modified_at": "string (ISO8601)"
    }
  ],
  "total": "number",
  "limit": "number",
  "offset": "number"
}
```

**Example**:

```bash
curl -X GET "http://localhost:8084/api/deduction-rules?limit=20&offset=0"
```

---

## Health Check

### Health Endpoint

Checks if the service is running.

**Endpoint**: `GET /health`

**Response** (200 OK):

```
Finance Service is healthy
```

---

## Error Responses

### Bad Request (400)

Returned when the request has invalid data.

```json
{
  "error": "Invalid request body: <details>"
}
```

### Not Found (404)

Returned when the requested resource doesn't exist.

```json
{
  "error": "Stipend not found"
}
```

### Server Error (500)

Returned when an internal server error occurs.

```json
{
  "error": "Failed to <operation>: <details>"
}
```

---

## Data Types

### Stipend Types

- `full-scholarship`: For scholarship students with limited deductions
- `self-funded`: For students funding themselves, with more deductions

### Payment Status

- `Pending`: Stipend created but not yet processed
- `Processed`: Stipend has been disbursed
- `Failed`: Stipend processing failed

### Processing Status

- `Pending`: Deduction awaiting review
- `Approved`: Deduction approved and ready to process
- `Processed`: Deduction has been applied
- `Rejected`: Deduction was rejected

### Payment Method

- `Bank_transfer`: Direct bank transfer
- `E-payment`: Electronic payment system
- Other custom methods as configured

---

## Rate Limiting

Currently, no rate limiting is implemented. Add rate limiting middleware as needed for production.

---

## CORS

Currently, CORS is not configured. Add CORS middleware if needed for cross-origin requests.

---

## API Examples

### Complete Stipend Processing Flow

```bash
# 1. Check available deduction rules
curl -X GET http://localhost:8084/api/deduction-rules

# 2. Calculate stipend with deductions (preview)
curl -X POST http://localhost:8084/api/stipends/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "stipend_type": "full-scholarship",
    "amount": 50000.00,
    "payment_method": "Bank_transfer",
    "journal_number": "JN-FS-001-2024"
  }'

# 3. Create stipend with calculated net amount
curl -X POST http://localhost:8084/api/stipends \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "stipend_type": "full-scholarship",
    "amount": 46500.00,
    "payment_method": "Bank_transfer",
    "journal_number": "JN-FS-001-2024",
    "notes": "Processed after deductions"
  }'

# 4. View created stipend
curl -X GET http://localhost:8084/api/stipends/f47ac10b-58cc-4372-a567-0e02b2c3d479

# 5. Get deductions applied
curl -X GET http://localhost:8084/api/stipends/f47ac10b-58cc-4372-a567-0e02b2c3d479/deductions

# 6. Mark as processed
curl -X PATCH http://localhost:8084/api/stipends/f47ac10b-58cc-4372-a567-0e02b2c3d479/payment-status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Processed",
    "payment_date": "2024-01-15T10:30:00Z"
  }'
```

---

## Money Transfer Endpoints

### 1. Initiate Transfer

Initiates a money transfer for a stipend to the student's bank account.

**Endpoint**: `POST /transfers/initiate`

**Request Body**:

```json
{
  "stipend_id": "string (UUID)",
  "payment_method": "string (BANK_TRANSFER|E_PAYMENT)"
}
```

**Response** (201 Created):

```json
{
  "id": "string (UUID - Transaction ID)",
  "stipend_id": "string (UUID)",
  "student_id": "string (UUID)",
  "amount": "number",
  "status": "string (PENDING|PROCESSING|SUCCESS|FAILED|CANCELLED)",
  "reference_number": "string (empty initially)",
  "error_message": "string (empty if success)",
  "payment_method": "string",
  "destination_account": "string",
  "destination_bank": "string",
  "initiated_at": "string (ISO8601)",
  "processed_at": "string (ISO8601, optional)",
  "completed_at": "string (ISO8601, optional)"
}
```

**Example**:

```bash
curl -X POST http://localhost:8084/api/transfers/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "stipend_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "payment_method": "BANK_TRANSFER"
  }'
```

---

### 2. Process Transfer

Processes a pending transfer by calling the payment gateway.

**Endpoint**: `POST /transfers/{transactionID}/process`

**Response** (200 OK):

```json
{
  "id": "string (UUID)",
  "status": "string (PROCESSING|SUCCESS|FAILED)",
  "reference_number": "string (populated on success)",
  "error_message": "string (populated on failure)",
  "completed_at": "string (ISO8601, populated on success)",
  ...
}
```

**Example**:

```bash
curl -X POST http://localhost:8084/api/transfers/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/process
```

---

### 3. Get Transfer Status

Retrieves the current status of a transfer.

**Endpoint**: `GET /transfers/{transactionID}/status`

**Response** (200 OK):

```json
{
  "id": "string (UUID)",
  "stipend_id": "string (UUID)",
  "student_id": "string (UUID)",
  "amount": "number",
  "status": "string",
  "reference_number": "string",
  "error_message": "string",
  "payment_method": "string",
  "destination_account": "string",
  "destination_bank": "string",
  "initiated_at": "string (ISO8601)",
  "processed_at": "string (ISO8601)",
  "completed_at": "string (ISO8601)"
}
```

**Example**:

```bash
curl -X GET http://localhost:8084/api/transfers/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/status
```

---

### 4. Get Transactions by Stipend

Retrieves all transactions (transfers) for a specific stipend.

**Endpoint**: `GET /stipends/{stipendID}/transactions`

**Response** (200 OK):

```json
[
  {
    "id": "string (UUID)",
    "stipend_id": "string (UUID)",
    "student_id": "string (UUID)",
    "amount": "number",
    "status": "string",
    ...
  }
]
```

**Example**:

```bash
curl -X GET http://localhost:8084/api/stipends/f47ac10b-58cc-4372-a567-0e02b2c3d479/transactions
```

---

### 5. Get Transactions by Student

Retrieves all transactions for a specific student.

**Endpoint**: `GET /students/{studentID}/transactions`

**Response** (200 OK):

```json
[
  {
    "id": "string (UUID)",
    "stipend_id": "string (UUID)",
    "student_id": "string (UUID)",
    "amount": "number",
    "status": "string",
    ...
  }
]
```

**Example**:

```bash
curl -X GET http://localhost:8084/api/students/550e8400-e29b-41d4-a716-446655440000/transactions
```

---

### 6. Cancel Transfer

Cancels a pending or processing transfer.

**Endpoint**: `POST /transfers/{transactionID}/cancel`

**Request Body**:

```json
{
  "reason": "string (cancellation reason)"
}
```

**Response** (200 OK):

```json
{
  "id": "string (UUID)",
  "status": "CANCELLED",
  "error_message": "string",
  ...
}
```

**Example**:

```bash
curl -X POST http://localhost:8084/api/transfers/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/cancel \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Student requested cancellation"
  }'
```

---

### 7. Retry Failed Transfer

Retries a failed transfer.

**Endpoint**: `POST /transfers/{transactionID}/retry`

**Response** (200 OK):

```json
{
  "id": "string (UUID)",
  "status": "string (PROCESSING|SUCCESS|FAILED)",
  "reference_number": "string (populated on success)",
  ...
}
```

**Example**:

```bash
curl -X POST http://localhost:8084/api/transfers/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/retry
```

---

## Transfer Status Codes

- `PENDING`: Transfer initiated but not yet processed
- `PROCESSING`: Transfer is being processed by payment gateway
- `SUCCESS`: Transfer completed successfully
- `FAILED`: Transfer failed during processing
- `CANCELLED`: Transfer was cancelled by user or system

---

## Complete Transfer Workflow Example

```bash
# 1. Calculate and create a stipend
curl -X POST http://localhost:8084/api/stipends/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "stipend_type": "full-scholarship",
    "amount": 50000.00
  }'

# 2. Create the stipend record with net amount
curl -X POST http://localhost:8084/api/stipends \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "stipend_type": "full-scholarship",
    "amount": 46500.00,
    "payment_method": "Bank_transfer",
    "journal_number": "JN-FS-001-2024"
  }'

# 3. Initiate the transfer
curl -X POST http://localhost:8084/api/transfers/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "stipend_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "payment_method": "BANK_TRANSFER"
  }'

# 4. Process the transfer (call payment gateway)
curl -X POST http://localhost:8084/api/transfers/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/process

# 5. Check transfer status
curl -X GET http://localhost:8084/api/transfers/a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6/status

# 6. View all student transactions
curl -X GET http://localhost:8084/api/students/550e8400-e29b-41d4-a716-446655440000/transactions
```

---

## Pagination

Pagination is supported on endpoints that return lists.

**Query Parameters**:

- `limit`: Number of items per page (default: varies by endpoint)
- `offset`: Number of items to skip from the beginning

**Response Format**:

```json
{
  "items": [...],
  "total": "number",
  "limit": "number",
  "offset": "number"
}
```

---

## Contact & Support

For issues or questions regarding the API, contact the Finance Service team.
