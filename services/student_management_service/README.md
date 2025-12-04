# Student Management Service

A comprehensive microservice for managing student information, academic programs, colleges, and stipend management in the RUB Student Portal system.

## Overview

This service handles all student-related operations including student registration, profile management, academic program management, college administration, and stipend eligibility and allocation tracking. It integrates with User Service for authentication and Banking Service for payment synchronization.

## Features

### Student Management

- **Student CRUD Operations**: Create, Read, Update, and Delete student records
- **Student Search**: Search students by name, email, or student ID
- **Filter by Program**: Get all students in a specific program
- **Filter by College**: Get all students in a specific college
- **Filter by Status**: Get students by their status (active, inactive, graduated, suspended)
- **Student ID Lookup**: Retrieve student information by their unique student ID
- **User Linking**: Students linked to User Service via UserID

### Program & College Management

- **Program Management**: CRUD operations for academic programs
- **College Management**: CRUD operations for colleges/institutions
- **Program-Stipend Linking**: Programs can be configured with stipend information
- **Hierarchical Structure**: Programs belong to colleges

### Stipend Management

- **Eligibility Checking**: Automatic calculation of stipend eligibility based on:
  - Student status (must be active)
  - Academic standing (must not be on probation/suspension)
  - GPA requirements (minimum 2.0)
  - Program stipend availability
- **Allocation Management**: Create and manage stipend allocations with approval workflow
- **Payment History**: Track all stipend payments and transactions
- **Status Tracking**: Monitor pending, approved, and disbursed allocations

### Reporting & Analytics

- **Student Summary Reports**: Aggregate statistics on students by college, program, and status
- **Stipend Statistics**: Financial reporting on allocations and disbursements
- **Custom Reports**: Filter reports by college or program
- **Eligibility Statistics**: Track eligible vs ineligible students

### Integration Features

- **Banking Service Sync**: Synchronize student bank details with Banking Service
- **Audit Logging**: Track all changes for compliance and security
- **Authentication Ready**: Middleware for JWT-based authentication (via API Gateway)

## Tech Stack

- **Language**: Go 1.22.1
- **Framework**: Chi Router v5
- **Database**: PostgreSQL
- **ORM**: GORM
- **Containerization**: Docker

## API Endpoints

### Student Management

| Method | Endpoint                               | Description                    |
| ------ | -------------------------------------- | ------------------------------ |
| GET    | `/api/students`                        | Get all students               |
| POST   | `/api/students`                        | Create a new student           |
| GET    | `/api/students/search?q={query}`       | Search students                |
| GET    | `/api/students/{id}`                   | Get student by database ID     |
| GET    | `/api/students/student-id/{studentId}` | Get student by student ID      |
| GET    | `/api/students/program/{programId}`    | Get students by program        |
| GET    | `/api/students/college/{collegeId}`    | Get students by college        |
| GET    | `/api/students/status/{status}`        | Get students by status         |
| PUT    | `/api/students/{id}`                   | Update a student               |
| DELETE | `/api/students/{id}`                   | Delete a student (soft delete) |

### Program Management

| Method | Endpoint             | Description          |
| ------ | -------------------- | -------------------- |
| GET    | `/api/programs`      | Get all programs     |
| POST   | `/api/programs`      | Create a new program |
| GET    | `/api/programs/{id}` | Get program by ID    |
| PUT    | `/api/programs/{id}` | Update a program     |

### College Management

| Method | Endpoint             | Description          |
| ------ | -------------------- | -------------------- |
| GET    | `/api/colleges`      | Get all colleges     |
| POST   | `/api/colleges`      | Create a new college |
| GET    | `/api/colleges/{id}` | Get college by ID    |
| PUT    | `/api/colleges/{id}` | Update a college     |

### Stipend Management

| Method | Endpoint                                   | Description                            |
| ------ | ------------------------------------------ | -------------------------------------- |
| GET    | `/api/stipend/eligibility/{studentId}`     | Check student stipend eligibility      |
| GET    | `/api/stipend/allocations`                 | Get all allocations (supports filters) |
| POST   | `/api/stipend/allocations`                 | Create new allocation                  |
| GET    | `/api/stipend/allocations/{id}`            | Get allocation by ID                   |
| PUT    | `/api/stipend/allocations/{id}`            | Update allocation (approve/reject)     |
| GET    | `/api/stipend/history`                     | Get all payment history                |
| POST   | `/api/stipend/history`                     | Record a payment                       |
| GET    | `/api/stipend/history/student/{studentId}` | Get student's payment history          |

### Reporting

| Method | Endpoint                                           | Description                       |
| ------ | -------------------------------------------------- | --------------------------------- |
| GET    | `/api/reports/students/summary`                    | Get student summary statistics    |
| GET    | `/api/reports/stipend/statistics`                  | Get stipend allocation statistics |
| GET    | `/api/reports/students/by-college?college_id={id}` | Get students by college           |
| GET    | `/api/reports/students/by-program?program_id={id}` | Get students by program           |

## Data Models

### Student Model

```json
{
  "user_id": 1,
  "student_id": "STU001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@student.rub.edu.bt",
  "phone_number": "+97517123456",
  "cid": "11234567890",
  "date_of_birth": "2000-01-15",
  "gender": "male",
  "permanent_address": "Thimphu, Bhutan",
  "current_address": "College Hostel",
  "program_id": 1,
  "college_id": 1,
  "year_of_study": 2,
  "semester": 4,
  "enrollment_date": "2022-08-01",
  "status": "active",
  "academic_standing": "good",
  "gpa": 3.5,
  "guardian_name": "Jane Doe",
  "guardian_phone_number": "+97517654321",
  "guardian_relation": "mother"
}
```

### Program Model

```json
{
  "code": "BSIT",
  "name": "Bachelor of Science in Information Technology",
  "description": "Four-year IT program",
  "level": "undergraduate",
  "duration_years": 4,
  "duration_semesters": 8,
  "college_id": 1,
  "has_stipend": true,
  "stipend_amount": 5000.0,
  "stipend_type": "semester",
  "is_active": true
}
```

### Stipend Allocation Model

```json
{
  "allocation_id": "STIP-2024-001",
  "student_id": 1,
  "amount": 5000.0,
  "allocation_date": "2024-01-15",
  "status": "approved",
  "approved_by": 10,
  "approval_date": "2024-01-16",
  "semester": 4,
  "academic_year": "2023-2024",
  "remarks": "Regular stipend allocation"
}
```

## Setup

### Prerequisites

- Go 1.22.1 or higher
- PostgreSQL database
- Docker (optional, for containerized deployment)

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
DATABASE_URL=postgresql://username:password@localhost:5432/student_portal_db
PORT=8084
BANKING_SERVICE_URL=http://localhost:8083
```

### Local Development

1. Install dependencies:

```bash
go mod download
```

2. Run the service:

```bash
go run main.go
```

The service will start on `http://localhost:8084`

### Docker Deployment

Build and run using Docker:

```bash
docker build -t student-management-service .
docker run -p 8084:8084 --env-file .env student-management-service
```

### Using Docker Compose

The service is configured in the main `docker-compose.yml` file:

```bash
docker-compose up student_management_service
```

## Database Schema

The service automatically creates the following tables:

- **students**: Main student information
- **colleges**: Academic colleges/institutions
- **programs**: Academic programs with stipend configuration
- **stipend_allocations**: Stipend allocation records
- **stipend_histories**: Payment transaction history
- **audit_logs**: Audit trail for all changes

All tables support soft deletes using GORM's `DeletedAt` field.

## Integration with Other Services

### User Service

- Students are linked via `user_id` field
- Authentication handled by API Gateway (JWT tokens)
- User info passed via headers (`X-User-ID`, `X-User-Role`)

### Banking Service

- Sync student bank details via REST API
- Endpoint: `POST /api/student-bank-details`
- Automatic retry on failures

### Financial Service (Future)

- Send approved stipend records for payment processing
- Receive payment confirmation callbacks

## Middleware & Security

### Authentication Middleware

- Validates user authentication via headers
- Extracts user ID and role from request context
- Returns 401 for unauthenticated requests

### Authorization

- `AdminOnly`: Restricts access to admin users
- `StudentOrAdmin`: Allows students or admins
- Context-based user identification

## User Story Implementation

This service implements the following user stories from `userstories.md`:

✅ **Epic: Student Profile Management**

- US1: Register New Student
- US2: View Student Profile
- US3: Update Student Profile
- US4: View Program Details

✅ **Epic: Stipend Management**

- US5: View Stipend History
- US7: Calculate Stipend Eligibility
- US8: Monitor Stipend Eligibility Statistics

✅ **Epic: Administrative Functions**

- US9: Deactivate or Archive Student
- US10: Assign College
- US11: Assign Academic Program
- US12: Generate Student Summary Report

✅ **Epic: System Integration**

- US13: Record Stipend Allocation
- US14: Sync Bank Details
- US15: Fetch User Identity
- US16: Send Stipend Record to Financial Service

## Development Notes

- Auto-migration creates/updates database schema on startup
- All timestamps managed automatically by GORM
- Soft deletes enabled for students, programs, colleges
- Connection pooling configured for optimal performance
- Search uses PostgreSQL's ILIKE for case-insensitive matching
- Audit logging tracks all data modifications

## Error Handling

The service returns appropriate HTTP status codes:

- `200 OK`: Successful GET/PUT requests
- `201 Created`: Successful POST requests
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing/invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Database or server errors

## Testing

Run tests:

```bash
go test ./... -v
```

Test coverage:

```bash
go test ./... -cover
```

## Future Enhancements

- [ ] Pagination for large datasets
- [ ] CSV/PDF export for reports
- [ ] Email notifications via Notification Service
- [ ] Advanced filtering and sorting
- [ ] Batch operations for data imports
- [ ] Student photo upload support
- [ ] Integration with Academic Records Service
- [ ] Real-time dashboards
- [ ] GraphQL API support
- [ ] Webhook support for external integrations

## License

This is part of the RUB Student Portal project.

## Overview

This service handles all student-related operations including student registration, profile management, and student information retrieval. It is part of the RUB Student Portal microservices architecture.

## Features

- **Student CRUD Operations**: Create, Read, Update, and Delete student records
- **Student Search**: Search students by name, email, or student ID
- **Filter by Program**: Get all students in a specific program
- **Filter by College**: Get all students in a specific college
- **Filter by Status**: Get students by their status (active, inactive, graduated, suspended)
- **Student ID Lookup**: Retrieve student information by their unique student ID

## Tech Stack

- **Language**: Go 1.22.1
- **Framework**: Chi Router v5
- **Database**: PostgreSQL
- **ORM**: GORM
- **Containerization**: Docker

## API Endpoints

### Student Management

| Method | Endpoint                               | Description                                   |
| ------ | -------------------------------------- | --------------------------------------------- |
| GET    | `/api/students`                        | Get all students                              |
| POST   | `/api/students`                        | Create a new student                          |
| GET    | `/api/students/search?q={query}`       | Search students by name, email, or student ID |
| GET    | `/api/students/{id}`                   | Get student by database ID                    |
| GET    | `/api/students/student-id/{studentId}` | Get student by student ID (e.g., STU001)      |
| GET    | `/api/students/program/{programId}`    | Get all students in a program                 |
| GET    | `/api/students/college/{collegeId}`    | Get all students in a college                 |
| GET    | `/api/students/status/{status}`        | Get students by status                        |
| PUT    | `/api/students/{id}`                   | Update a student                              |
| DELETE | `/api/students/{id}`                   | Delete a student (soft delete)                |

## Student Model

```json
{
  "student_id": "STU001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@student.rub.edu.bt",
  "phone_number": "+97517123456",
  "date_of_birth": "2000-01-15",
  "gender": "male",
  "permanent_address": "Thimphu, Bhutan",
  "current_address": "College Hostel, Phuntsholing",
  "program_id": 1,
  "college_id": 1,
  "year_of_study": 2,
  "semester": 4,
  "enrollment_date": "2022-08-01",
  "graduation_date": "",
  "status": "active",
  "guardian_name": "Jane Doe",
  "guardian_phone_number": "+97517654321",
  "guardian_relation": "mother"
}
```

## Setup

### Prerequisites

- Go 1.22.1 or higher
- PostgreSQL database
- Docker (optional, for containerized deployment)

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
DATABASE_URL=postgresql://username:password@localhost:5432/student_portal_db
PORT=8084
```

### Local Development

1. Install dependencies:

```bash
go mod download
```

2. Run the service:

```bash
go run main.go
```

The service will start on `http://localhost:8084`

### Docker Deployment

Build and run using Docker:

```bash
docker build -t student-management-service .
docker run -p 8084:8084 --env-file .env student-management-service
```

### Using Docker Compose

The service is configured in the main `docker-compose.yml` file:

```bash
docker-compose up student_management_service
```

## Database Schema

The service automatically creates the following table structure:

- **students**: Main student information table with fields for personal details, academic information, and guardian information
- Supports soft deletes using GORM's `DeletedAt` field
- Unique constraints on `student_id` and `email`

## Integration with Other Services

This service is designed to work alongside:

- **User Services**: For authentication and user management
- **Banking Services**: For student financial information
- **API Gateway**: Routes requests to this service

## Development Notes

- The service uses GORM's Auto-Migration to create/update database schema
- All timestamps are automatically managed by GORM
- Soft deletes are enabled for student records
- Connection pooling is configured for optimal database performance
- Search functionality uses PostgreSQL's ILIKE for case-insensitive matching

## Error Handling

The service returns appropriate HTTP status codes:

- `200 OK`: Successful GET/PUT/DELETE requests
- `201 Created`: Successful POST requests
- `400 Bad Request`: Invalid input data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Database or server errors

## Future Enhancements

- Pagination for large datasets
- Advanced filtering options
- Batch operations for student imports
- Integration with academic records service
- Student photo upload support
- Email notifications for important events
- Export functionality (CSV, PDF)

## License

This is part of the RUB Student Portal project.
