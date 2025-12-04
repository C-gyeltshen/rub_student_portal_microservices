# Postman Testing Guide - Student Management Service

## Service Information

- **Base URL**: `http://localhost:8084`
- **Content-Type**: `application/json` (for POST/PUT requests)
- **Service Status**: HTTP Server on port 8084

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [College Endpoints](#college-endpoints)
3. [Program Endpoints](#program-endpoints)
4. [Student Endpoints](#student-endpoints)
5. [Stipend Endpoints](#stipend-endpoints)
6. [Report Endpoints](#report-endpoints)
7. [Testing Sequence](#testing-sequence)

---

## Prerequisites

### 1. Start the Service

```bash
cd services/student_management_service
go run main.go
```

**Expected output:**

```
Database connected and all models migrated successfully.
Student Management Service starting on :8084
gRPC server listening on :50054
```

### 2. Postman Setup

1. Open Postman
2. Create a new Collection: "Student Management Service"
3. Set environment variable: `base_url` = `http://localhost:8084`

---

## College Endpoints

### 1. Create College (POST)

**Endpoint:** `POST http://localhost:8084/api/colleges`

**Headers:**

```
Content-Type: application/json
```

**Request Body:**

```json
{
  "code": "CST",
  "name": "College of Science & Technology",
  "description": "Engineering and Technology programs",
  "location": "Phuentsholing",
  "is_active": true
}
```

**Required Fields:**

- `code` (string) - Unique college code
- `name` (string) - College name

**Optional Fields:**

- `description` (string)
- `location` (string)
- `is_active` (boolean) - defaults to true

**Expected Response (201 Created):**

```json
{
  "ID": 1,
  "CreatedAt": "2025-11-26T15:30:00Z",
  "UpdatedAt": "2025-11-26T15:30:00Z",
  "DeletedAt": null,
  "code": "CST",
  "name": "College of Science & Technology",
  "description": "Engineering and Technology programs",
  "location": "Phuentsholing",
  "is_active": true
}
```

---

### 2. Get All Colleges (GET)

**Endpoint:** `GET http://localhost:8084/api/colleges`

**No body required**

**Expected Response (200 OK):**

```json
[
  {
    "ID": 1,
    "code": "CST",
    "name": "College of Science & Technology",
    "location": "Phuentsholing",
    "is_active": true
  }
]
```

---

### 3. Get College by ID (GET)

**Endpoint:** `GET http://localhost:8084/api/colleges/1`

**No body required**

---

### 4. Update College (PUT)

**Endpoint:** `PUT http://localhost:8084/api/colleges/1`

**Request Body (partial update):**

```json
{
  "location": "Rinchending, Phuentsholing",
  "description": "Updated description"
}
```

---

## Program Endpoints

### 1. Create Program (POST)

**Endpoint:** `POST http://localhost:8084/api/programs`

**Headers:**

```
Content-Type: application/json
```

**Request Body:**

```json
{
  "code": "BSIT",
  "name": "Bachelor of Science in Information Technology",
  "description": "4-year undergraduate IT program",
  "level": "undergraduate",
  "duration_years": 4,
  "duration_semesters": 8,
  "college_id": 1,
  "has_stipend": true,
  "stipend_amount": 5000.0,
  "stipend_type": "monthly",
  "is_active": true
}
```

**Required Fields:**

- `code` (string) - Unique program code
- `name` (string) - Program name

**Optional Fields:**

- `description` (string)
- `level` (string) - "undergraduate" or "postgraduate"
- `duration_years` (number)
- `duration_semesters` (number)
- `college_id` (number) - Must exist in colleges table
- `has_stipend` (boolean) - Default false
- `stipend_amount` (number) - Required if has_stipend is true
- `stipend_type` (string) - "monthly", "semester", or "annual"
- `is_active` (boolean) - Default true

**Expected Response (201 Created):**

```json
{
  "ID": 1,
  "code": "BSIT",
  "name": "Bachelor of Science in Information Technology",
  "level": "undergraduate",
  "duration_years": 4,
  "duration_semesters": 8,
  "college_id": 1,
  "college": {
    "ID": 1,
    "code": "CST",
    "name": "College of Science & Technology"
  },
  "has_stipend": true,
  "stipend_amount": 5000,
  "stipend_type": "monthly"
}
```

---

### 2. Get All Programs (GET)

**Endpoint:** `GET http://localhost:8084/api/programs`

**No body required**

**Expected Response (200 OK):**

```json
[
  {
    "ID": 1,
    "code": "BSIT",
    "name": "Bachelor of Science in Information Technology",
    "college": {
      "ID": 1,
      "name": "College of Science & Technology"
    },
    "has_stipend": true,
    "stipend_amount": 5000
  }
]
```

---

### 3. Get Program by ID (GET)

**Endpoint:** `GET http://localhost:8084/api/programs/1`

---

### 4. Update Program (PUT)

**Endpoint:** `PUT http://localhost:8084/api/programs/1`

**Request Body:**

```json
{
  "stipend_amount": 6000.0,
  "duration_years": 4
}
```

---

## Student Endpoints

### 1. Create Student (POST)

**Endpoint:** `POST http://localhost:8084/api/students`

**Headers:**

```
Content-Type: application/json
```

**Request Body - Minimum Required:**

```json
{
  "user_id": 1,
  "student_id": "02220123",
  "first_name": "Tshering",
  "last_name": "Wangpo",
  "email": "tshering.wangpo@rub.edu.bt"
}
```

**Request Body - Complete:**

```json
{
  "user_id": 1,
  "student_id": "02220123",
  "first_name": "Tshering",
  "last_name": "Wangpo",
  "email": "tshering.wangpo@rub.edu.bt",
  "phone_number": "+97517123456",
  "date_of_birth": "2003-05-15",
  "gender": "Male",
  "cid": "11234567890",
  "permanent_address": "Thimphu, Bhutan",
  "current_address": "Phuentsholing, Bhutan",
  "program_id": 1,
  "college_id": 1,
  "year_of_study": 2,
  "semester": 3,
  "enrollment_date": "2022-08-01",
  "graduation_date": "2026-06-30",
  "status": "active",
  "academic_standing": "good",
  "gpa": 3.5,
  "guardian_name": "Dorji Wangchuk",
  "guardian_phone_number": "+97517654321",
  "guardian_relation": "Father"
}
```

**Field Reference:**

| Field                   | Type   | Required | Example                 | Notes                                          |
| ----------------------- | ------ | -------- | ----------------------- | ---------------------------------------------- |
| `user_id`               | number | ✅       | `1`                     | Links to User Service (must exist)             |
| `student_id`            | string | ✅       | `"02220123"`            | RUB student ID (unique)                        |
| `first_name`            | string | ✅       | `"Tshering"`            | Student's first name                           |
| `last_name`             | string | ✅       | `"Wangpo"`              | Student's last name                            |
| `email`                 | string | ✅       | `"tshering@rub.edu.bt"` | Must be unique                                 |
| `phone_number`          | string | ❌       | `"+97517123456"`        | Contact number                                 |
| `date_of_birth`         | string | ❌       | `"2003-05-15"`          | Format: YYYY-MM-DD                             |
| `gender`                | string | ❌       | `"Male"` or `"Female"`  | Gender                                         |
| `cid`                   | string | ❌       | `"11234567890"`         | Citizenship ID (unique if provided)            |
| `permanent_address`     | string | ❌       | `"Thimphu, Bhutan"`     | Home address                                   |
| `current_address`       | string | ❌       | `"Phuentsholing"`       | Current residential address                    |
| `program_id`            | number | ❌       | `1`                     | Must exist in programs table                   |
| `college_id`            | number | ❌       | `1`                     | Must exist in colleges table                   |
| `year_of_study`         | number | ❌       | `2`                     | Current academic year (1-4)                    |
| `semester`              | number | ❌       | `3`                     | Current semester (1-8)                         |
| `enrollment_date`       | string | ❌       | `"2022-08-01"`          | Admission date                                 |
| `graduation_date`       | string | ❌       | `"2026-06-30"`          | Expected graduation                            |
| `status`                | string | ❌       | `"active"`              | `active`, `inactive`, `graduated`, `suspended` |
| `academic_standing`     | string | ❌       | `"good"`                | `good`, `probation`                            |
| `gpa`                   | number | ❌       | `3.5`                   | Grade point average (0.0-4.0)                  |
| `guardian_name`         | string | ❌       | `"Dorji Wangchuk"`      | Guardian's full name                           |
| `guardian_phone_number` | string | ❌       | `"+97517654321"`        | Guardian's contact                             |
| `guardian_relation`     | string | ❌       | `"Father"`              | Relationship to student                        |

**❗ Do NOT Include:**

- `ID` - Auto-generated
- `CreatedAt` - Auto-generated
- `UpdatedAt` - Auto-generated
- `DeletedAt` - For soft deletes
- `program` - Loaded from database
- `college` - Loaded from database

**Expected Response (201 Created):**

```json
{
  "ID": 1,
  "CreatedAt": "2025-11-26T16:00:00Z",
  "UpdatedAt": "2025-11-26T16:00:00Z",
  "user_id": 1,
  "student_id": "02220123",
  "first_name": "Tshering",
  "last_name": "Wangpo",
  "email": "tshering.wangpo@rub.edu.bt",
  "phone_number": "+97517123456",
  "program_id": 1,
  "college_id": 1,
  "status": "active",
  "gpa": 3.5,
  "academic_standing": "good"
}
```

---

### 2. Get All Students (GET)

**Endpoint:** `GET http://localhost:8084/api/students`

**No body required**

**Expected Response (200 OK):**

```json
[
  {
    "ID": 1,
    "student_id": "02220123",
    "first_name": "Tshering",
    "last_name": "Wangpo",
    "email": "tshering.wangpo@rub.edu.bt",
    "status": "active",
    "gpa": 3.5
  }
]
```

---

### 3. Get Student by Database ID (GET)

**Endpoint:** `GET http://localhost:8084/api/students/1`

**No body required**

---

### 4. Get Student by Student ID (GET)

**Endpoint:** `GET http://localhost:8084/api/students/student-id/02220123`

**No body required**

**Use Case:** Search by RUB student ID instead of database ID

---

### 5. Search Students (GET)

**Endpoint:** `GET http://localhost:8084/api/students/search?q=Tshering`

**Query Parameters:**

- `q` (required) - Search query

**Search Examples:**

```
/api/students/search?q=Tshering    # Search by first name
/api/students/search?q=Wangpo      # Search by last name
/api/students/search?q=rub.edu.bt  # Search by email
/api/students/search?q=02220123    # Search by student ID
```

**Features:**

- Case-insensitive search
- Searches across: first_name, last_name, email, student_id
- Partial matching

---

### 6. Get Students by Program (GET)

**Endpoint:** `GET http://localhost:8084/api/students/program/1`

**Returns all students enrolled in program ID 1**

---

### 7. Get Students by College (GET)

**Endpoint:** `GET http://localhost:8084/api/students/college/1`

**Returns all students in college ID 1**

---

### 8. Get Students by Status (GET)

**Endpoint:** `GET http://localhost:8084/api/students/status/active`

**Valid status values:**

- `active`
- `inactive`
- `graduated`
- `suspended`

---

### 9. Update Student (PUT)

**Endpoint:** `PUT http://localhost:8084/api/students/1`

**Headers:**

```
Content-Type: application/json
```

**Request Body (partial update):**

```json
{
  "gpa": 3.7,
  "semester": 4,
  "academic_standing": "good"
}
```

**You can update any field except:**

- `ID`
- `CreatedAt`
- `user_id` (should not change)
- `student_id` (should not change)

---

### 10. Delete Student (DELETE)

**Endpoint:** `DELETE http://localhost:8084/api/students/1`

**No body required**

**Expected Response (200 OK):**

```json
{
  "message": "Student deleted successfully"
}
```

**Note:** This is a soft delete - sets `deleted_at` timestamp

---

## Stipend Endpoints

### 1. Check Stipend Eligibility (GET)

**Endpoint:** `GET http://localhost:8084/api/stipend/eligibility/1`

**No body required**

**Expected Response (200 OK) - Eligible:**

```json
{
  "student_id": 1,
  "is_eligible": true,
  "reasons": ["All eligibility criteria met"],
  "expected_amount": 5000.0,
  "academic_standing": "good",
  "attendance_rate": 0,
  "has_pending_issues": false
}
```

**Expected Response - Not Eligible:**

```json
{
  "student_id": 1,
  "is_eligible": false,
  "reasons": ["GPA below minimum requirement"],
  "expected_amount": 0,
  "academic_standing": "probation"
}
```

**Eligibility Rules:**

1. Student status must be "active"
2. Academic standing must NOT be "probation" or "suspended"
3. GPA must be >= 2.0
4. Program must have `has_stipend: true`

---

### 2. Create Stipend Allocation (POST)

**Endpoint:** `POST http://localhost:8084/api/stipend/allocations`

**Headers:**

```
Content-Type: application/json
```

**Request Body:**

```json
{
  "allocation_id": "STIP-2025-001",
  "student_id": 1,
  "amount": 5000.0,
  "allocation_date": "2025-01-01",
  "status": "pending",
  "semester": 3,
  "academic_year": "2024-2025",
  "remarks": "Monthly stipend for semester 3"
}
```

**Required Fields:**

- `allocation_id` (string) - Unique identifier
- `student_id` (number) - Must exist and be eligible
- `amount` (number)

**Optional Fields:**

- `allocation_date` (string)
- `status` (string) - Default "pending"
- `approved_by` (number) - User ID who approved
- `approval_date` (string)
- `semester` (number)
- `academic_year` (string)
- `remarks` (string)

**Status Values:**

- `pending` - Initial state
- `approved` - Authorized by admin
- `rejected` - Denied
- `disbursed` - Payment sent

**Expected Response (201 Created):**

```json
{
  "ID": 1,
  "allocation_id": "STIP-2025-001",
  "student_id": 1,
  "amount": 5000,
  "status": "pending",
  "semester": 3,
  "academic_year": "2024-2025"
}
```

---

### 3. Get All Stipend Allocations (GET)

**Endpoint:** `GET http://localhost:8084/api/stipend/allocations`

**No body required**

**With Query Filters:**

```
/api/stipend/allocations?status=pending
/api/stipend/allocations?student_id=1
/api/stipend/allocations?status=approved&student_id=1
```

---

### 4. Get Stipend Allocation by ID (GET)

**Endpoint:** `GET http://localhost:8084/api/stipend/allocations/1`

---

### 5. Update Stipend Allocation (PUT)

**Endpoint:** `PUT http://localhost:8084/api/stipend/allocations/1`

**Request Body (Approve allocation):**

```json
{
  "status": "approved",
  "approved_by": 5,
  "approval_date": "2025-11-26",
  "remarks": "Approved by Finance Admin"
}
```

**Request Body (Reject allocation):**

```json
{
  "status": "rejected",
  "remarks": "Student does not meet GPA requirement"
}
```

---

### 6. Get Stipend Payment History (GET)

**Endpoint:** `GET http://localhost:8084/api/stipend/history`

**With Query Filters:**

```
/api/stipend/history?student_id=1
/api/stipend/history?status=success
```

---

### 7. Get Student Stipend History (GET)

**Endpoint:** `GET http://localhost:8084/api/stipend/history/student/1`

**Returns all payment history for student ID 1**

---

### 8. Create Stipend Payment Record (POST)

**Endpoint:** `POST http://localhost:8084/api/stipend/history`

**Request Body:**

```json
{
  "transaction_id": "TXN-2025-001",
  "student_id": 1,
  "allocation_id": 1,
  "amount": 5000.0,
  "payment_date": "2025-01-15",
  "transaction_status": "success",
  "payment_method": "bank_transfer",
  "bank_reference": "BNK123456789",
  "remarks": "January stipend payment"
}
```

**Required Fields:**

- `transaction_id` (string) - Unique
- `student_id` (number)
- `amount` (number)

**Transaction Status Values:**

- `success`
- `failed`
- `pending`

**Payment Method Values:**

- `bank_transfer`
- `cash`
- `mobile_payment`

---

## Report Endpoints

### 1. Student Summary Report (GET)

**Endpoint:** `GET http://localhost:8084/api/reports/students/summary`

**No body required**

**Expected Response (200 OK):**

```json
{
  "total_students": 150,
  "active_students": 120,
  "inactive_students": 20,
  "graduated_students": 10,
  "by_college": {
    "College of Science & Technology": 80,
    "College of Natural Resources": 70
  },
  "by_program": {
    "Bachelor of Science in IT": 45,
    "Bachelor of Education": 35
  },
  "by_status": {
    "active": 120,
    "inactive": 20,
    "graduated": 10
  },
  "stipend_eligible": 100,
  "stipend_ineligible": 50
}
```

---

### 2. Stipend Statistics Report (GET)

**Endpoint:** `GET http://localhost:8084/api/reports/stipend/statistics`

**No body required**

**Expected Response (200 OK):**

```json
{
  "total_allocations": 250,
  "total_amount": 1250000.0,
  "pending_allocations": 50,
  "approved_allocations": 180,
  "disbursed_amount": 900000.0,
  "by_program": {
    "Bachelor of Science in IT": 600000.0,
    "Bachelor of Education": 400000.0
  },
  "by_college": {
    "College of Science & Technology": 750000.0,
    "College of Natural Resources": 500000.0
  }
}
```

---

### 3. Students by College Report (GET)

**Endpoint:** `GET http://localhost:8084/api/reports/students/by-college?college_id=1`

**Query Parameters:**

- `college_id` (required) - College ID

**Returns detailed list of all students in the specified college**

---

### 4. Students by Program Report (GET)

**Endpoint:** `GET http://localhost:8084/api/reports/students/by-program?program_id=1`

**Query Parameters:**

- `program_id` (required) - Program ID

**Returns detailed list of all students in the specified program**

---

## Testing Sequence

### Quick Test Flow (Recommended Order)

Follow this sequence to test all functionality:

#### Step 1: Create Foundation Data

```
1. POST /api/colleges          → Create "CST" college
2. GET  /api/colleges          → Verify college created
3. POST /api/programs          → Create "BSIT" program (use college_id: 1)
4. GET  /api/programs          → Verify program created
```

#### Step 2: Create Students

```
5. POST /api/students          → Create student 1 (minimal fields)
6. POST /api/students          → Create student 2 (complete fields)
7. GET  /api/students          → List all students
8. GET  /api/students/1        → Get student by ID
9. GET  /api/students/student-id/02220123  → Get by student ID
```

#### Step 3: Test Search & Filters

```
10. GET /api/students/search?q=Tshering    → Search students
11. GET /api/students/program/1            → Filter by program
12. GET /api/students/college/1            → Filter by college
13. GET /api/students/status/active        → Filter by status
```

#### Step 4: Test Stipend Management

```
14. GET  /api/stipend/eligibility/1        → Check eligibility
15. POST /api/stipend/allocations          → Create allocation
16. GET  /api/stipend/allocations          → List allocations
17. PUT  /api/stipend/allocations/1        → Approve allocation
18. POST /api/stipend/history              → Record payment
19. GET  /api/stipend/history/student/1    → View student history
```

#### Step 5: Generate Reports

```
20. GET /api/reports/students/summary           → Student statistics
21. GET /api/reports/stipend/statistics         → Stipend analytics
22. GET /api/reports/students/by-college?college_id=1
23. GET /api/reports/students/by-program?program_id=1
```

#### Step 6: Test Updates & Deletes

```
24. PUT    /api/students/1     → Update student GPA
25. PUT    /api/programs/1     → Update program stipend
26. DELETE /api/students/1     → Soft delete student
```

---

## Common Error Responses

### 400 Bad Request

```json
{
  "error": "Missing required fields"
}
```

**Fix:** Check that all required fields are included

### 404 Not Found

```json
{
  "error": "Student not found"
}
```

**Fix:** Verify the ID exists in database

### 500 Internal Server Error

```json
{
  "error": "pq: duplicate key value violates unique constraint \"idx_students_email\""
}
```

**Fix:** Email already exists, use a different email

---

## Sample Data Sets

### Sample College

```json
{
  "code": "CNR",
  "name": "College of Natural Resources",
  "location": "Lobesa",
  "is_active": true
}
```

### Sample Program

```json
{
  "code": "BED",
  "name": "Bachelor of Education",
  "level": "undergraduate",
  "duration_years": 4,
  "college_id": 1,
  "has_stipend": true,
  "stipend_amount": 4500.0,
  "stipend_type": "monthly"
}
```

### Sample Students

```json
{
  "user_id": 2,
  "student_id": "02220124",
  "first_name": "Karma",
  "last_name": "Dorji",
  "email": "karma.dorji@rub.edu.bt",
  "program_id": 1,
  "college_id": 1,
  "gpa": 3.2,
  "status": "active",
  "academic_standing": "good"
}
```

```json
{
  "user_id": 3,
  "student_id": "02220125",
  "first_name": "Pema",
  "last_name": "Choden",
  "email": "pema.choden@rub.edu.bt",
  "program_id": 1,
  "college_id": 1,
  "gpa": 2.8,
  "status": "active",
  "academic_standing": "good"
}
```

---

## Tips for Postman

### 1. Environment Variables

Create environment variables for reusability:

```
base_url = http://localhost:8084
college_id = 1
program_id = 1
student_id = 1
```

Use in requests: `{{base_url}}/api/students/{{student_id}}`

### 2. Save Responses

After creating a college/program, save the ID from the response to use in subsequent requests.

### 3. Collections

Organize requests into folders:

- Colleges
- Programs
- Students
- Stipend
- Reports

### 4. Tests Tab

Add test scripts to automatically extract IDs:

```javascript
var jsonData = pm.response.json();
pm.environment.set("college_id", jsonData.ID);
```

### 5. Pre-request Scripts

Set timestamps automatically:

```javascript
pm.environment.set("current_date", new Date().toISOString().split("T")[0]);
```

---

## Troubleshooting

### Service Not Running

**Error:** `Could not get response - Error: connect ECONNREFUSED 127.0.0.1:8084`

**Solution:**

```bash
cd services/student_management_service
go run main.go
```

### Database Not Connected

**Error:** `DATABASE_URL environment variable is not set`

**Solution:** Create `.env` file with:

```
DATABASE_URL=postgresql://rubadmin:rubpassword@localhost:5432/student_service_db
PORT=8084
```

### Foreign Key Errors

**Error:** `foreign key constraint fails`

**Solution:** Create college and program before creating students

### Unique Constraint Violations

**Error:** `duplicate key value violates unique constraint`

**Solution:** Use different values for:

- `student_id`
- `email`
- `cid`
- `user_id`

---

**Last Updated:** November 26, 2025  
**Service Version:** 1.0.0
