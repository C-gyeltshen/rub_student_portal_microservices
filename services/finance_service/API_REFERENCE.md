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
