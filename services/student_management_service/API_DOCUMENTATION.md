# Student Management Service - API Documentation

## Base URL

```
http://localhost:8084
```

## Authentication

All endpoints expect authentication headers from the API Gateway:

- `X-User-ID`: The authenticated user's ID
- `X-User-Role`: The user's role (student, admin, finance_officer)

## Response Format

All responses are in JSON format with appropriate HTTP status codes.

---

## Student Endpoints

### Create Student

**POST** `/api/students`

**Request Body:**

```json
{
  "user_id": 1,
  "student_id": "STU2024001",
  "first_name": "Tshering",
  "last_name": "Wangpo",
  "email": "tshering.wangpo@student.rub.edu.bt",
  "phone_number": "+97517123456",
  "cid": "11234567890",
  "program_id": 1,
  "college_id": 1,
  "year_of_study": 1,
  "semester": 1,
  "status": "active"
}
```

**Response:** `201 Created`

```json
{
  "id": 1,
  "user_id": 1,
  "student_id": "STU2024001",
  "first_name": "Tshering",
  ...
  "created_at": "2024-11-25T00:00:00Z"
}
```

### Get All Students

**GET** `/api/students`

**Response:** `200 OK`

```json
[
  {
    "id": 1,
    "student_id": "STU2024001",
    "first_name": "Tshering",
    ...
  }
]
```

### Search Students

**GET** `/api/students/search?q={query}`

**Query Parameters:**

- `q` (required): Search term

**Response:** `200 OK`

### Get Student by ID

**GET** `/api/students/{id}`

**Response:** `200 OK`

### Update Student

**PUT** `/api/students/{id}`

**Request Body:**

```json
{
  "phone_number": "+97517999888",
  "email": "newemail@student.rub.edu.bt",
  "status": "active"
}
```

**Response:** `200 OK`

---

## Program Endpoints

### Create Program

**POST** `/api/programs`

**Request Body:**

```json
{
  "code": "BSIT",
  "name": "Bachelor of Science in Information Technology",
  "description": "4-year IT program",
  "level": "undergraduate",
  "duration_years": 4,
  "duration_semesters": 8,
  "college_id": 1,
  "has_stipend": true,
  "stipend_amount": 5000.0,
  "stipend_type": "semester"
}
```

**Response:** `201 Created`

### Get All Programs

**GET** `/api/programs`

**Response:** `200 OK`

```json
[
  {
    "id": 1,
    "code": "BSIT",
    "name": "Bachelor of Science in Information Technology",
    "has_stipend": true,
    "stipend_amount": 5000.0,
    "college": {
      "id": 1,
      "name": "College of Science and Technology"
    }
  }
]
```

---

## College Endpoints

### Create College

**POST** `/api/colleges`

**Request Body:**

```json
{
  "code": "CST",
  "name": "College of Science and Technology",
  "description": "Premier technology institute",
  "location": "Phuntsholing",
  "is_active": true
}
```

**Response:** `201 Created`

### Get All Colleges

**GET** `/api/colleges`

**Response:** `200 OK`

---

## Stipend Endpoints

### Check Eligibility

**GET** `/api/stipend/eligibility/{studentId}`

**Response:** `200 OK`

```json
{
  "student_id": 1,
  "is_eligible": true,
  "reasons": ["All eligibility criteria met"],
  "expected_amount": 5000.0,
  "academic_standing": "good",
  "attendance_rate": 95.5,
  "has_pending_issues": false
}
```

### Create Allocation

**POST** `/api/stipend/allocations`

**Request Body:**

```json
{
  "allocation_id": "STIP-2024-001",
  "student_id": 1,
  "amount": 5000.0,
  "allocation_date": "2024-01-15",
  "semester": 1,
  "academic_year": "2023-2024",
  "status": "pending"
}
```

**Response:** `201 Created`

### Get Allocations

**GET** `/api/stipend/allocations?status={status}&student_id={id}`

**Query Parameters:**

- `status` (optional): Filter by status (pending, approved, rejected, disbursed)
- `student_id` (optional): Filter by student

**Response:** `200 OK`

### Approve/Reject Allocation

**PUT** `/api/stipend/allocations/{id}`

**Request Body:**

```json
{
  "status": "approved",
  "approved_by": 10,
  "approval_date": "2024-01-16",
  "remarks": "Approved for disbursement"
}
```

**Response:** `200 OK`

### Get Stipend History

**GET** `/api/stipend/history?student_id={id}&status={status}`

**Query Parameters:**

- `student_id` (optional): Filter by student
- `status` (optional): Filter by transaction status

**Response:** `200 OK`

```json
[
  {
    "id": 1,
    "transaction_id": "TXN-2024-001",
    "student_id": 1,
    "amount": 5000.0,
    "payment_date": "2024-01-20",
    "transaction_status": "success",
    "bank_reference": "BANK123456"
  }
]
```

### Get Student Payment History

**GET** `/api/stipend/history/student/{studentId}`

**Response:** `200 OK`

---

## Report Endpoints

### Student Summary

**GET** `/api/reports/students/summary`

**Response:** `200 OK`

```json
{
  "total_students": 150,
  "active_students": 120,
  "inactive_students": 20,
  "graduated_students": 10,
  "by_college": {
    "College of Science and Technology": 80,
    "College of Natural Resources": 70
  },
  "by_program": {
    "BSc IT": 50,
    "BE Civil": 30
  },
  "stipend_eligible": 100,
  "stipend_ineligible": 50
}
```

### Stipend Statistics

**GET** `/api/reports/stipend/statistics`

**Response:** `200 OK`

```json
{
  "total_allocations": 100,
  "total_amount": 500000.0,
  "pending_allocations": 20,
  "approved_allocations": 80,
  "disbursed_amount": 400000.0,
  "by_program": {
    "BSc IT": 250000.0,
    "BE Civil": 150000.0
  },
  "by_college": {
    "CST": 300000.0,
    "CNR": 200000.0
  }
}
```

### Students by College Report

**GET** `/api/reports/students/by-college?college_id={id}`

**Response:** `200 OK`

### Students by Program Report

**GET** `/api/reports/students/by-program?program_id={id}`

**Response:** `200 OK`

---

## Error Responses

### 400 Bad Request

```json
{
  "error": "Missing required fields: student_id, first_name, email"
}
```

### 401 Unauthorized

```json
{
  "error": "Unauthorized: Missing user authentication"
}
```

### 403 Forbidden

```json
{
  "error": "Forbidden: Admin access required"
}
```

### 404 Not Found

```json
{
  "error": "Student not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "Database connection failed"
}
```

---

## Common Use Cases

### 1. Register New Student

```bash
curl -X POST http://localhost:8084/api/students \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -H "X-User-Role: admin" \
  -d '{
    "user_id": 1,
    "student_id": "STU2024001",
    "first_name": "Tshering",
    "last_name": "Wangpo",
    "email": "tshering@student.rub.edu.bt",
    "program_id": 1,
    "college_id": 1
  }'
```

### 2. Check Stipend Eligibility

```bash
curl http://localhost:8084/api/stipend/eligibility/1 \
  -H "X-User-ID: 1" \
  -H "X-User-Role: student"
```

### 3. Generate Student Report

```bash
curl http://localhost:8084/api/reports/students/summary \
  -H "X-User-ID: 10" \
  -H "X-User-Role: admin"
```

### 4. View Stipend History

```bash
curl http://localhost:8084/api/stipend/history/student/1 \
  -H "X-User-ID: 1" \
  -H "X-User-Role: student"
```

---

## Postman Collection

Import this collection into Postman for easy testing:
[Link to Postman collection would go here]

## Rate Limiting

Currently no rate limiting is implemented. This will be handled at the API Gateway level.

## Versioning

Current API version: `v1` (implicit in all endpoints)
Future versions will use `/api/v2/` prefix.
