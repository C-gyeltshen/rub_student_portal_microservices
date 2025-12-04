# Bulk Student Creation API

This document provides comprehensive information about the bulk student creation endpoints in the Student Management Service.

## Table of Contents

- [Overview](#overview)
- [Endpoints](#endpoints)
- [Bulk JSON Creation](#bulk-json-creation)
- [Bulk CSV Creation](#bulk-csv-creation)
- [Testing Examples](#testing-examples)

---

## Overview

The Student Management Service provides two endpoints for bulk student creation:

1. **POST /api/students/bulk** - Create multiple students from JSON array
2. **POST /api/students/bulk/csv** - Create multiple students from CSV file upload

Both endpoints validate data, provide detailed error reporting, and use database transactions to ensure data integrity.

---

## Endpoints

### 1. Bulk JSON Creation

**Endpoint**: `POST /api/students/bulk`  
**Content-Type**: `application/json`  
**Description**: Creates multiple students from a JSON array

#### Required Fields

- `student_id` (string, unique)
- `first_name` (string)
- `last_name` (string)
- `email` (string, unique)

#### Optional Fields

- `phone_number` (string)
- `cid` (string, unique) - Citizenship ID
- `date_of_birth` (string)
- `gender` (string)
- `program_id` (uint)
- `college_id` (uint)
- `year_of_study` (int)
- `semester` (int)
- `gpa` (float64)
- `status` (string) - default: "active"
- `enrollment_date` (string)
- `graduation_date` (string)
- `academic_standing` (string)
- `permanent_address` (string)
- `current_address` (string)
- `guardian_name` (string)
- `guardian_phone_number` (string)
- `guardian_relation` (string)

#### Example Request

```bash
curl -X POST http://localhost:8086/api/students/bulk \
  -H "Content-Type: application/json" \
  -d '[
    {
      "student_id": "02230001",
      "first_name": "Tenzin",
      "last_name": "Dorji",
      "email": "02230001.gcit@rub.edu.bt",
      "phone_number": "+97517123456",
      "cid": "11234567890",
      "date_of_birth": "2004-03-15",
      "gender": "Male",
      "program_id": 1,
      "college_id": 1,
      "year_of_study": 2,
      "semester": 4,
      "gpa": 3.45,
      "status": "active",
      "enrollment_date": "2023-01-15",
      "permanent_address": "Thimphu, Bhutan",
      "current_address": "Phuentsholing, Bhutan",
      "academic_standing": "good",
      "guardian_name": "Sonam Dorji",
      "guardian_phone_number": "+97517654321",
      "guardian_relation": "Father"
    },
    {
      "student_id": "02230002",
      "first_name": "Pema",
      "last_name": "Wangmo",
      "email": "02230002.gcit@rub.edu.bt",
      "phone_number": "+97517234567",
      "cid": "11234567891",
      "date_of_birth": "2003-08-22",
      "gender": "Female",
      "program_id": 2,
      "college_id": 1,
      "year_of_study": 3,
      "semester": 5,
      "gpa": 3.78,
      "status": "active",
      "enrollment_date": "2022-01-10",
      "permanent_address": "Paro, Bhutan",
      "current_address": "Phuentsholing, Bhutan",
      "academic_standing": "good",
      "guardian_name": "Karma Wangmo",
      "guardian_phone_number": "+97517345678",
      "guardian_relation": "Mother"
    },
    {
      "student_id": "02230003",
      "first_name": "Karma",
      "last_name": "Tshering",
      "email": "02230003.gcit@rub.edu.bt",
      "phone_number": "+97517345679",
      "cid": "11234567892",
      "date_of_birth": "2004-11-08",
      "gender": "Male",
      "program_id": 1,
      "college_id": 1,
      "year_of_study": 2,
      "semester": 3,
      "gpa": 3.12,
      "status": "active",
      "enrollment_date": "2023-01-15"
    }
  ]'
```

#### Success Response

**Status Code**: `201 Created`

```json
{
  "success": true,
  "created_count": 3,
  "total_count": 3,
  "errors": [],
  "message": "Bulk creation completed"
}
```

#### Validation Error Response

**Status Code**: `400 Bad Request`

```json
{
  "success": false,
  "message": "Validation failed for 2 student(s)",
  "validation_errors": [
    {
      "index": 1,
      "missing_fields": ["email"],
      "received_data": {
        "student_id": "02230004",
        "first_name": "Sonam",
        "last_name": "Choden",
        "email": "",
        "cid": "11234567893",
        "gender": "Female"
      }
    },
    {
      "index": 3,
      "missing_fields": ["first_name", "last_name"],
      "received_data": {
        "student_id": "02230005",
        "first_name": "",
        "last_name": "",
        "email": "02230005.gcit@rub.edu.bt",
        "cid": "11234567894",
        "gender": "Male"
      }
    }
  ],
  "total_students": 4
}
```

#### Partial Success Response

**Status Code**: `201 Created`

```json
{
  "success": true,
  "created_count": 2,
  "total_count": 3,
  "errors": [
    "Failed to create student 02230006: duplicate key value violates unique constraint \"students_email_key\""
  ],
  "message": "Bulk creation completed"
}
```

---

### 2. Bulk CSV Creation

**Endpoint**: `POST /api/students/bulk/csv`  
**Content-Type**: `multipart/form-data`  
**Description**: Creates multiple students from a CSV file upload

#### CSV Format

The CSV file must have a header row with the following column names (in any order):

**Required Columns**:

- `student_id`
- `first_name`
- `last_name`
- `email`

**Optional Columns**:

- `phone_number`
- `cid`
- `date_of_birth`
- `gender`
- `program_id`
- `college_id`
- `year_of_study`
- `semester`
- `gpa`
- `status`
- `enrollment_date`
- `graduation_date`
- `academic_standing`
- `permanent_address`
- `current_address`
- `guardian_name`
- `guardian_phone_number`
- `guardian_relation`

#### Example CSV File

```csv
student_id,first_name,last_name,email,phone_number,cid,date_of_birth,gender,program_id,college_id,year_of_study,semester,gpa,status,enrollment_date,permanent_address,current_address,academic_standing,guardian_name,guardian_phone_number,guardian_relation
02230001,Tenzin,Dorji,02230001.gcit@rub.edu.bt,+97517123456,11234567890,2004-03-15,Male,1,1,2,4,3.45,active,2023-01-15,"Thimphu, Bhutan","Phuentsholing, Bhutan",good,Sonam Dorji,+97517654321,Father
02230002,Pema,Wangmo,02230002.gcit@rub.edu.bt,+97517234567,11234567891,2003-08-22,Female,2,1,3,5,3.78,active,2022-01-10,"Paro, Bhutan","Phuentsholing, Bhutan",good,Karma Wangmo,+97517345678,Mother
02230003,Karma,Tshering,02230003.gcit@rub.edu.bt,+97517345679,11234567892,2004-11-08,Male,1,1,2,3,3.12,active,2023-01-15,"Punakha, Bhutan","Phuentsholing, Bhutan",good,Ugyen Tshering,+97517456789,Father
02230004,Deki,Zangmo,02230004.gcit@rub.edu.bt,+97517456790,11234567893,2004-05-20,Female,3,1,1,2,3.65,active,2024-01-10,"Bumthang, Bhutan","Phuentsholing, Bhutan",good,Namgay Zangmo,+97517567890,Mother
```

#### Example Request with cURL

```bash
curl -X POST http://localhost:8086/api/students/bulk/csv \
  -F "file=@students.csv"
```

#### Example Request with Postman

1. Set method to `POST`
2. Set URL to `http://localhost:8086/api/students/bulk/csv`
3. Go to **Body** tab
4. Select **form-data**
5. Add key `file` with type `File`
6. Choose your CSV file

#### Success Response

**Status Code**: `201 Created`

```json
{
  "success": true,
  "created_count": 4,
  "total_count": 4,
  "errors": [],
  "parse_errors": [],
  "message": "CSV bulk creation completed"
}
```

#### Validation Error Response

**Status Code**: `400 Bad Request`

```json
{
  "success": false,
  "message": "Validation failed for 2 student(s)",
  "validation_errors": [
    {
      "row": 3,
      "missing_fields": ["email"],
      "received_data": {
        "student_id": "02230005",
        "first_name": "Sonam",
        "last_name": "Choden",
        "email": ""
      }
    },
    {
      "row": 5,
      "missing_fields": ["first_name", "last_name"],
      "received_data": {
        "student_id": "02230006",
        "first_name": "",
        "last_name": "",
        "email": "02230006.gcit@rub.edu.bt"
      }
    }
  ],
  "parse_errors": [],
  "total_rows": 5
}
```

---

## Testing Examples

### Using cURL (JSON)

**Minimal Required Fields**:

```bash
curl -X POST http://localhost:8086/api/students/bulk \
  -H "Content-Type: application/json" \
  -d '[
    {
      "student_id": "02230010",
      "first_name": "Test",
      "last_name": "Student",
      "email": "test.student@rub.edu.bt"
    }
  ]'
```

**With All Fields**:

```bash
curl -X POST http://localhost:8086/api/students/bulk \
  -H "Content-Type: application/json" \
  -d '[
    {
      "student_id": "02230011",
      "first_name": "Complete",
      "last_name": "Example",
      "email": "complete.example@rub.edu.bt",
      "phone_number": "+97517111111",
      "cid": "11111111111",
      "date_of_birth": "2004-01-01",
      "gender": "Male",
      "program_id": 1,
      "college_id": 1,
      "year_of_study": 1,
      "semester": 1,
      "gpa": 3.50,
      "status": "active",
      "enrollment_date": "2024-01-10",
      "graduation_date": "2027-12-15",
      "permanent_address": "Thimphu, Bhutan",
      "current_address": "Phuentsholing, Bhutan",
      "academic_standing": "good",
      "guardian_name": "Guardian Name",
      "guardian_phone_number": "+97517222222",
      "guardian_relation": "Father"
    }
  ]'
```

### Using cURL (CSV)

```bash
curl -X POST http://localhost:8086/api/students/bulk/csv \
  -F "file=@/path/to/students.csv"
```

### Using JavaScript/Fetch API (JSON)

```javascript
const students = [
  {
    student_id: "02230020",
    first_name: "JavaScript",
    last_name: "Example",
    email: "js.example@rub.edu.bt",
    phone_number: "+97517333333",
    program_id: 1,
    college_id: 1,
    year_of_study: 2,
    semester: 3,
    gpa: 3.25,
    status: "active",
  },
];

fetch("http://localhost:8086/api/students/bulk", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify(students),
})
  .then((response) => response.json())
  .then((data) => console.log("Success:", data))
  .catch((error) => console.error("Error:", error));
```

### Using JavaScript/Fetch API (CSV)

```javascript
const fileInput = document.querySelector('input[type="file"]');
const formData = new FormData();
formData.append("file", fileInput.files[0]);

fetch("http://localhost:8086/api/students/bulk/csv", {
  method: "POST",
  body: formData,
})
  .then((response) => response.json())
  .then((data) => console.log("Success:", data))
  .catch((error) => console.error("Error:", error));
```

---

## Error Handling

### Common Errors

1. **Missing Required Fields**

   - Returns validation errors with specific field names
   - No records are created if validation fails

2. **Duplicate Values**

   - `student_id`, `email`, and `cid` must be unique
   - Returns database constraint violation error

3. **Invalid JSON Format**

   - Returns parse error if JSON is malformed

4. **Invalid CSV Format**

   - Returns parse errors for malformed CSV rows
   - Continues processing valid rows

5. **File Upload Errors**
   - File size limit: 10 MB
   - Only CSV files should be uploaded

### Transaction Behavior

- All operations are wrapped in database transactions
- If all students fail validation, no records are created
- If some students are created successfully, those are committed
- Partial success is reported with detailed error messages

---

## Best Practices

1. **Validate Data Client-Side**: Check required fields before sending requests
2. **Handle Errors Gracefully**: Parse error responses and display user-friendly messages
3. **Use Appropriate Endpoint**: Use CSV for large datasets, JSON for programmatic access
4. **Check Response**: Always check `success` field and `created_count` vs `total_count`
5. **Unique Constraints**: Ensure `student_id`, `email`, and `cid` are unique across all students
6. **Date Formats**: Use consistent date format (YYYY-MM-DD recommended)
7. **Test Small Batches**: Test with small datasets before processing large files

---

## Performance Considerations

- **Batch Size**: Recommended maximum of 1000 students per request
- **CSV vs JSON**: CSV is more efficient for very large datasets
- **Database Transactions**: All students are created in a single transaction
- **Validation**: All validation happens before any database writes

---

## Reference Data

### Available Colleges

Before creating students, you need to create colleges first. Use the following endpoints:

**Get All Colleges**: `GET /api/colleges`  
**Create College**: `POST /api/colleges`

#### College Structure

```json
{
  "code": "string (unique)",
  "name": "string (required)",
  "description": "string",
  "location": "string",
  "is_active": "boolean (default: true)"
}
```

#### Example RUB Colleges

```bash
# Create College of Science and Technology (CST)
curl -X POST http://localhost:8086/api/colleges \
  -H "Content-Type: application/json" \
  -d '{
    "code": "CST",
    "name": "College of Science and Technology",
    "description": "College of Science and Technology, Phuentsholing",
    "location": "Phuentsholing",
    "is_active": true
  }'

# Create Gedu College of Business Studies (GCBS)
curl -X POST http://localhost:8086/api/colleges \
  -H "Content-Type: application/json" \
  -d '{
    "code": "GCBS",
    "name": "Gedu College of Business Studies",
    "description": "Gedu College of Business Studies",
    "location": "Gedu",
    "is_active": true
  }'

# Create College of Natural Resources (CNR)
curl -X POST http://localhost:8086/api/colleges \
  -H "Content-Type: application/json" \
  -d '{
    "code": "CNR",
    "name": "College of Natural Resources",
    "description": "College of Natural Resources, Lobesa",
    "location": "Lobesa, Punakha",
    "is_active": true
  }'

# Create Paro College of Education (PCE)
curl -X POST http://localhost:8086/api/colleges \
  -H "Content-Type: application/json" \
  -d '{
    "code": "PCE",
    "name": "Paro College of Education",
    "description": "Paro College of Education",
    "location": "Paro",
    "is_active": true
  }'

# Create Sherubtse College
curl -X POST http://localhost:8086/api/colleges \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SC",
    "name": "Sherubtse College",
    "description": "Sherubtse College, Kanglung",
    "location": "Kanglung, Trashigang",
    "is_active": true
  }'

# Create Jigme Namgyel Engineering College (JNEC)
curl -X POST http://localhost:8086/api/colleges \
  -H "Content-Type: application/json" \
  -d '{
    "code": "JNEC",
    "name": "Jigme Namgyel Engineering College",
    "description": "Jigme Namgyel Engineering College, Dewathang",
    "location": "Dewathang, Samdrup Jongkhar",
    "is_active": true
  }'
```

---

### Available Programs

Before creating students, you need to create programs first. Use the following endpoints:

**Get All Programs**: `GET /api/programs`  
**Create Program**: `POST /api/programs`

#### Program Structure

```json
{
  "code": "string (unique)",
  "name": "string (required)",
  "description": "string",
  "level": "string (undergraduate/postgraduate)",
  "duration_years": "int",
  "duration_semesters": "int",
  "college_id": "uint (required)",
  "has_stipend": "boolean (default: false)",
  "stipend_amount": "float64",
  "stipend_type": "string (monthly/semester/annual)",
  "is_active": "boolean (default: true)"
}
```

#### Example Programs

**Note**: Replace `college_id` with the actual ID returned after creating colleges.

```bash
# BSc in Information Technology (CST - Assuming college_id = 1)
curl -X POST http://localhost:8086/api/programs \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BSIT",
    "name": "Bachelor of Science in Information Technology",
    "description": "Undergraduate program in Information Technology",
    "level": "undergraduate",
    "duration_years": 4,
    "duration_semesters": 8,
    "college_id": 1,
    "has_stipend": true,
    "stipend_amount": 3000.00,
    "stipend_type": "monthly",
    "is_active": true
  }'

# BE in Electronics and Communication Engineering (CST)
curl -X POST http://localhost:8086/api/programs \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BEECE",
    "name": "Bachelor of Engineering in Electronics and Communication",
    "description": "Undergraduate program in Electronics and Communication Engineering",
    "level": "undergraduate",
    "duration_years": 4,
    "duration_semesters": 8,
    "college_id": 1,
    "has_stipend": true,
    "stipend_amount": 3000.00,
    "stipend_type": "monthly",
    "is_active": true
  }'

# BBA - Bachelor of Business Administration (GCBS - Assuming college_id = 2)
curl -X POST http://localhost:8086/api/programs \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BBA",
    "name": "Bachelor of Business Administration",
    "description": "Undergraduate program in Business Administration",
    "level": "undergraduate",
    "duration_years": 4,
    "duration_semesters": 8,
    "college_id": 2,
    "has_stipend": true,
    "stipend_amount": 3000.00,
    "stipend_type": "monthly",
    "is_active": true
  }'

# BEd - Bachelor of Education (PCE - Assuming college_id = 4)
curl -X POST http://localhost:8086/api/programs \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BED",
    "name": "Bachelor of Education",
    "description": "Undergraduate program in Education",
    "level": "undergraduate",
    "duration_years": 4,
    "duration_semesters": 8,
    "college_id": 4,
    "has_stipend": true,
    "stipend_amount": 3000.00,
    "stipend_type": "monthly",
    "is_active": true
  }'

# BSc in Environmental Science (CNR - Assuming college_id = 3)
curl -X POST http://localhost:8086/api/programs \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BSES",
    "name": "Bachelor of Science in Environmental Science",
    "description": "Undergraduate program in Environmental Science",
    "level": "undergraduate",
    "duration_years": 4,
    "duration_semesters": 8,
    "college_id": 3,
    "has_stipend": true,
    "stipend_amount": 3000.00,
    "stipend_type": "monthly",
    "is_active": true
  }'

# BE in Civil Engineering (JNEC - Assuming college_id = 6)
curl -X POST http://localhost:8086/api/programs \
  -H "Content-Type: application/json" \
  -d '{
    "code": "BECE",
    "name": "Bachelor of Engineering in Civil Engineering",
    "description": "Undergraduate program in Civil Engineering",
    "level": "undergraduate",
    "duration_years": 4,
    "duration_semesters": 8,
    "college_id": 6,
    "has_stipend": true,
    "stipend_amount": 3000.00,
    "stipend_type": "monthly",
    "is_active": true
  }'
```

#### Retrieve College and Program IDs

After creating colleges and programs, retrieve their IDs to use in student creation:

```bash
# Get all colleges with their IDs
curl http://localhost:8086/api/colleges

# Get all programs with their IDs
curl http://localhost:8086/api/programs

# Get specific college by ID
curl http://localhost:8086/api/colleges/1

# Get specific program by ID
curl http://localhost:8086/api/programs/1
```

**Example Response**:

```json
[
  {
    "ID": 1,
    "code": "CST",
    "name": "College of Science and Technology",
    "description": "College of Science and Technology, Phuentsholing",
    "location": "Phuentsholing",
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

### Important Notes

1. **Create Colleges First**: Always create colleges before creating programs
2. **Create Programs Second**: Create programs with valid `college_id` references
3. **Create Students Last**: Create students with valid `college_id` and `program_id` references
4. **Foreign Key Constraints**: The database enforces referential integrity
5. **Use Actual IDs**: Replace example college_id values with actual IDs from your database

---

## Service Information

- **Service Port**: 8086 (configurable via PORT environment variable)
- **Database**: PostgreSQL (via GORM)
- **Router**: Chi v5
- **Max File Size**: 10 MB for CSV uploads
