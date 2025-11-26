# Student Management Service - Complete Technical Documentation

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Project Structure](#project-structure)
4. [Database Models](#database-models)
5. [API Endpoints](#api-endpoints)
6. [gRPC Services](#grpc-services)
7. [File & Folder Documentation](#file--folder-documentation)
8. [Technologies Used](#technologies-used)

---

## Overview

The **Student Management Service** is a microservice responsible for managing all student-related operations in the RUB Student Portal system. It handles student enrollment, academic records, program management, college information, stipend allocation, and comprehensive reporting.

### Key Features

- ✅ Student CRUD operations with comprehensive data management
- ✅ College and Program management
- ✅ Stipend eligibility checking and allocation tracking
- ✅ Payment history tracking
- ✅ Advanced search and filtering capabilities
- ✅ Comprehensive reporting and analytics
- ✅ Hybrid HTTP REST + gRPC architecture
- ✅ Service-to-service communication with User and Banking services
- ✅ Audit logging for compliance
- ✅ Soft delete support for data recovery

### Service Information

- **HTTP Port**: 8084
- **gRPC Port**: 50054
- **Database**: PostgreSQL (student_service_db)
- **Framework**: Go with Chi router
- **ORM**: GORM

---

## Architecture

### Hybrid Communication Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Applications                       │
│              (Postman, React Frontend, Mobile)               │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP REST (Port 8084)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│           Student Management Service (Port 8084)            │
│                                                               │
│  ┌──────────────────┐          ┌─────────────────────┐      │
│  │  HTTP REST API   │          │   gRPC Server       │      │
│  │  (Chi Router)    │          │   (Port 50054)      │      │
│  │  - 28 Endpoints  │          │   - 10 RPC Methods  │      │
│  └──────────────────┘          └─────────────────────┘      │
│           │                              │                   │
│           └──────────┬───────────────────┘                   │
│                      ▼                                       │
│           ┌─────────────────────┐                           │
│           │  Business Logic     │                           │
│           │  (Handlers)         │                           │
│           └──────────┬──────────┘                           │
│                      ▼                                       │
│           ┌─────────────────────┐                           │
│           │   Database Layer    │                           │
│           │   (GORM ORM)        │                           │
│           └──────────┬──────────┘                           │
└──────────────────────┼──────────────────────────────────────┘
                       │
                       ▼
           ┌───────────────────────┐
           │   PostgreSQL 15       │
           │ (student_service_db)  │
           └───────────────────────┘

         gRPC Service-to-Service Communication
                       │
       ┌───────────────┴───────────────┐
       ▼                               ▼
┌─────────────────┐          ┌──────────────────┐
│  User Service   │          │ Banking Service  │
│  (Port 50052)   │          │  (Port 50053)    │
│  - GetUser      │          │  - GetBankDetails│
│  - ValidateToken│          │  - VerifyAccount │
│  - GetUserRole  │          │  - UpsertDetails │
└─────────────────┘          └──────────────────┘
```

### Data Flow

1. **External Requests** → HTTP REST API (Port 8084)
2. **Internal Requests** → gRPC API (Port 50054)
3. **Service Communication** → gRPC Clients (User & Banking Services)
4. **Data Persistence** → PostgreSQL via GORM

---

## Project Structure

```
student_management_service/
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── .env                         # Environment variables (DATABASE_URL, ports)
├── Dockerfile                   # Docker containerization config
│
├── database/                    # Database connection & configuration
│   └── db.go                    # GORM setup, connection pooling, AutoMigrate
│
├── models/                      # Data models (structs)
│   ├── student.go               # Student entity (25+ fields)
│   ├── program.go               # Program & College entities
│   ├── stipend.go               # Stipend allocations & history
│   └── audit.go                 # Audit logging model
│
├── handlers/                    # HTTP request handlers
│   ├── student_handler.go       # Student CRUD operations
│   ├── program_handler.go       # Program & College operations
│   ├── stipend_handler.go       # Stipend management
│   └── report_handler.go        # Analytics & reporting
│
├── grpc/                        # gRPC implementation
│   ├── server/
│   │   └── student_server.go    # gRPC server (10 RPC methods)
│   └── client/
│       ├── user_client.go       # User Service gRPC client
│       └── banking_client.go    # Banking Service gRPC client
│
├── proto/                       # Protocol Buffer definitions
│   ├── student/
│   │   └── student.proto        # Student Service API contract
│   ├── user/
│   │   └── user.proto           # User Service API contract
│   └── banking/
│       └── banking.proto        # Banking Service API contract
│
├── pb/                          # Generated protobuf code
│   ├── student/
│   │   ├── student.pb.go        # Message definitions
│   │   └── student_grpc.pb.go   # gRPC service definitions
│   ├── user/
│   └── banking/
│
├── middleware/                  # HTTP middleware
│   └── auth.go                  # Authentication helpers
│
├── utils/                       # Utility functions
│   └── audit.go                 # Audit logging utilities
│
└── Documentation/
    ├── README.md                # Service overview
    ├── API_DOCUMENTATION.md     # HTTP REST API docs
    ├── GRPC_DOCUMENTATION.md    # gRPC implementation guide
    ├── GRPC_SUMMARY.md          # gRPC quick reference
    └── userstories.md           # User stories & requirements
```

---

## Database Models

### 1. Student Model (`models/student.go`)

**Purpose**: Core entity representing a student with personal, academic, and guardian information.

**Database Table**: `students`

**Key Features**:

- Unique student ID (RUB student identifier)
- Links to User Service via `user_id`
- Foreign keys to Program and College
- Soft delete support via `DeletedAt`
- Comprehensive address and guardian info

**Code Structure**:

```go
type Student struct {
    gorm.Model                    // ID, CreatedAt, UpdatedAt, DeletedAt

    // Link to User Service
    UserID      uint             // Foreign key to User Service

    // Student Information
    StudentID   string           // RUB Student ID (unique)
    FirstName   string           // Required
    LastName    string           // Required
    Email       string           // Unique, required
    PhoneNumber string
    DateOfBirth string
    Gender      string
    CID         string           // Citizenship ID (unique)

    // Address
    PermanentAddress string
    CurrentAddress   string

    // Academic Information
    ProgramID        uint        // Foreign key to programs table
    Program          Program     // Belongs to relationship
    CollegeID        uint        // Foreign key to colleges table
    College          College     // Belongs to relationship
    YearOfStudy      int
    Semester         int
    EnrollmentDate   string
    GraduationDate   string
    Status           string      // active, inactive, graduated, suspended
    AcademicStanding string      // good, probation
    GPA              float64

    // Guardian Information
    GuardianName        string
    GuardianPhoneNumber string
    GuardianRelation    string
}
```

**Database Indexes**:

- `user_id` (unique index)
- `student_id` (unique index)
- `email` (unique index)
- `cid` (unique index)
- `deleted_at` (for soft delete queries)

---

### 2. College Model (`models/program.go`)

**Purpose**: Represents academic colleges/institutions within RUB.

**Database Table**: `colleges`

**Code Structure**:

```go
type College struct {
    gorm.Model
    Code        string    // Unique college code (e.g., "CNR", "CST")
    Name        string    // Full name (e.g., "College of Natural Resources")
    Description string
    Location    string
    IsActive    bool      // Enable/disable colleges
}
```

**Examples**:

- Code: "CST", Name: "College of Science & Technology"
- Code: "CNR", Name: "College of Natural Resources"

---

### 3. Program Model (`models/program.go`)

**Purpose**: Represents academic programs offered by colleges.

**Database Table**: `programs`

**Key Features**:

- Linked to parent College
- Stipend configuration (amount, type, eligibility)
- Duration tracking (years and semesters)

**Code Structure**:

```go
type Program struct {
    gorm.Model
    Code             string     // Unique program code (e.g., "BSIT", "BED")
    Name             string     // Full name
    Description      string
    Level            string     // undergraduate, postgraduate
    DurationYears    int
    DurationSemesters int
    CollegeID        uint       // Foreign key to colleges
    College          College    // Belongs to relationship

    // Stipend Configuration
    HasStipend       bool       // Does this program offer stipend?
    StipendAmount    float64    // Amount per payment
    StipendType      string     // monthly, semester, annual

    IsActive         bool
}
```

**Examples**:

```json
{
  "code": "BSIT",
  "name": "Bachelor of Science in Information Technology",
  "level": "undergraduate",
  "duration_years": 4,
  "duration_semesters": 8,
  "has_stipend": true,
  "stipend_amount": 5000.0,
  "stipend_type": "monthly"
}
```

---

### 4. StipendAllocation Model (`models/stipend.go`)

**Purpose**: Tracks stipend allocations to students with approval workflow.

**Database Table**: `stipend_allocations`

**Key Features**:

- Links to student via `student_id`
- Approval workflow (pending → approved → disbursed)
- Tracks approver and approval date
- Semester and academic year tracking

**Code Structure**:

```go
type StipendAllocation struct {
    gorm.Model
    AllocationID   string     // Unique allocation identifier
    StudentID      uint       // Foreign key to students
    Amount         float64
    AllocationDate string
    Status         string     // pending, approved, rejected, disbursed
    ApprovedBy     uint       // User ID who approved
    ApprovalDate   string
    Semester       int
    AcademicYear   string
    Remarks        string
}
```

**Workflow**:

1. **pending**: Initial allocation created
2. **approved**: Authorized by admin
3. **disbursed**: Payment sent to student
4. **rejected**: Allocation denied

---

### 5. StipendHistory Model (`models/stipend.go`)

**Purpose**: Immutable record of stipend payments for audit trail.

**Database Table**: `stipend_histories`

**Code Structure**:

```go
type StipendHistory struct {
    gorm.Model
    TransactionID     string     // Unique transaction ID
    StudentID         uint       // Foreign key to students
    AllocationID      uint       // Links to allocation
    Amount            float64
    PaymentDate       string
    TransactionStatus string     // success, failed, pending
    PaymentMethod     string     // bank_transfer, cash
    BankReference     string     // External bank reference
    Remarks           string
}
```

---

### 6. AuditLog Model (`models/audit.go`)

**Purpose**: Records all changes to critical data for compliance and debugging.

**Database Table**: `audit_logs`

**Code Structure**:

```go
type AuditLog struct {
    gorm.Model
    EntityType   string        // student, stipend, program
    EntityID     uint          // ID of the changed entity
    Action       string        // create, update, delete
    UserID       uint          // Who made the change
    Changes      string        // JSON of what changed (JSONB type)
    IPAddress    string        // Request IP
    UserAgent    string        // Browser/client info
    Timestamp    time.Time
}
```

**Example Entry**:

```json
{
  "entity_type": "student",
  "entity_id": 123,
  "action": "update",
  "user_id": 45,
  "changes": "{\"gpa\": {\"old\": 3.2, \"new\": 3.5}}",
  "ip_address": "192.168.1.100",
  "timestamp": "2025-11-26T15:30:00Z"
}
```

---

## API Endpoints

### Student Endpoints (10 endpoints)

| Method | Endpoint                               | Handler                   | Purpose                 |
| ------ | -------------------------------------- | ------------------------- | ----------------------- |
| GET    | `/api/students`                        | `GetStudents()`           | List all students       |
| POST   | `/api/students`                        | `CreateStudent()`         | Create new student      |
| GET    | `/api/students/{id}`                   | `GetStudentById()`        | Get by database ID      |
| GET    | `/api/students/student-id/{studentId}` | `GetStudentByStudentId()` | Get by RUB student ID   |
| GET    | `/api/students/search?q={query}`       | `SearchStudents()`        | Search by name/email/ID |
| GET    | `/api/students/program/{programId}`    | `GetStudentsByProgram()`  | Filter by program       |
| GET    | `/api/students/college/{collegeId}`    | `GetStudentsByCollege()`  | Filter by college       |
| GET    | `/api/students/status/{status}`        | `GetStudentsByStatus()`   | Filter by status        |
| PUT    | `/api/students/{id}`                   | `UpdateStudent()`         | Update student          |
| DELETE | `/api/students/{id}`                   | `DeleteStudent()`         | Soft delete student     |

### Program Endpoints (4 endpoints)

| Method | Endpoint             | Handler            | Purpose            |
| ------ | -------------------- | ------------------ | ------------------ |
| GET    | `/api/programs`      | `GetPrograms()`    | List all programs  |
| POST   | `/api/programs`      | `CreateProgram()`  | Create new program |
| GET    | `/api/programs/{id}` | `GetProgramById()` | Get program by ID  |
| PUT    | `/api/programs/{id}` | `UpdateProgram()`  | Update program     |

### College Endpoints (4 endpoints)

| Method | Endpoint             | Handler            | Purpose            |
| ------ | -------------------- | ------------------ | ------------------ |
| GET    | `/api/colleges`      | `GetColleges()`    | List all colleges  |
| POST   | `/api/colleges`      | `CreateCollege()`  | Create new college |
| GET    | `/api/colleges/{id}` | `GetCollegeById()` | Get college by ID  |
| PUT    | `/api/colleges/{id}` | `UpdateCollege()`  | Update college     |

### Stipend Endpoints (7 endpoints)

| Method | Endpoint                                   | Handler                      | Purpose              |
| ------ | ------------------------------------------ | ---------------------------- | -------------------- |
| GET    | `/api/stipend/eligibility/{studentId}`     | `CheckStipendEligibility()`  | Check eligibility    |
| GET    | `/api/stipend/allocations`                 | `GetStipendAllocations()`    | List allocations     |
| POST   | `/api/stipend/allocations`                 | `CreateStipendAllocation()`  | Create allocation    |
| GET    | `/api/stipend/allocations/{id}`            | `GetStipendAllocationById()` | Get allocation       |
| PUT    | `/api/stipend/allocations/{id}`            | `UpdateStipendAllocation()`  | Update allocation    |
| GET    | `/api/stipend/history`                     | `GetStipendHistory()`        | List payment history |
| GET    | `/api/stipend/history/student/{studentId}` | `GetStudentStipendHistory()` | Student's history    |

### Report Endpoints (4 endpoints)

| Method | Endpoint                           | Handler                        | Purpose            |
| ------ | ---------------------------------- | ------------------------------ | ------------------ |
| GET    | `/api/reports/students/summary`    | `GenerateStudentSummary()`     | Student statistics |
| GET    | `/api/reports/stipend/statistics`  | `GenerateStipendStatistics()`  | Stipend analytics  |
| GET    | `/api/reports/students/by-college` | `GetStudentsByCollegeReport()` | College report     |
| GET    | `/api/reports/students/by-program` | `GetStudentsByProgramReport()` | Program report     |

**Total: 29 HTTP REST Endpoints**

---

## gRPC Services

### Student Service (10 RPC Methods)

**Server**: `grpc/server/student_server.go`  
**Proto**: `proto/student/student.proto`

| RPC Method                | Request                        | Response                     | Purpose        |
| ------------------------- | ------------------------------ | ---------------------------- | -------------- |
| `GetStudent`              | `GetStudentRequest`            | `StudentResponse`            | Get by ID      |
| `GetStudentByStudentId`   | `GetStudentByStudentIdRequest` | `StudentResponse`            | Get by RUB ID  |
| `CreateStudent`           | `CreateStudentRequest`         | `StudentResponse`            | Create student |
| `UpdateStudent`           | `UpdateStudentRequest`         | `StudentResponse`            | Update student |
| `DeleteStudent`           | `DeleteStudentRequest`         | `DeleteResponse`             | Delete student |
| `ListStudents`            | `ListStudentsRequest`          | `ListStudentsResponse`       | List all       |
| `SearchStudents`          | `SearchStudentsRequest`        | `ListStudentsResponse`       | Search         |
| `GetStudentsByProgram`    | `GetByProgramRequest`          | `ListStudentsResponse`       | By program     |
| `GetStudentsByCollege`    | `GetByCollegeRequest`          | `ListStudentsResponse`       | By college     |
| `CheckStipendEligibility` | `StipendEligibilityRequest`    | `StipendEligibilityResponse` | Eligibility    |

---

## File & Folder Documentation

### 1. `main.go` - Application Entry Point (126 lines)

**Purpose**: Bootstrap application, start HTTP and gRPC servers concurrently

**Key Functions**:

#### `main()`

```go
func main() {
    // 1. Load environment variables from .env file
    godotenv.Load()

    // 2. Connect to PostgreSQL database
    database.Connect()

    // 3. Start gRPC server in background goroutine
    go startGRPCServer()

    // 4. Start HTTP server on main thread (blocking)
    startHTTPServer()
}
```

**Why Goroutine?**

- gRPC server runs concurrently in background
- HTTP server blocks main thread
- Both servers operate simultaneously on different ports

#### `startGRPCServer()` - Port 50054

```go
func startGRPCServer() {
    grpcPort := os.Getenv("GRPC_PORT")  // 50054

    lis, err := net.Listen("tcp", ":"+grpcPort)
    grpcServer := grpc.NewServer()

    // Register StudentService implementation
    pb.RegisterStudentServiceServer(grpcServer, &grpcserver.StudentServer{})

    grpcServer.Serve(lis)  // Blocking call
}
```

#### `startHTTPServer()` - Port 8084

```go
func startHTTPServer() {
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.Logger)     // Log all HTTP requests
    r.Use(middleware.Recoverer)  // Recover from panics

    // Register 29 HTTP routes
    // Student routes (10)
    r.Get("/api/students", handlers.GetStudents)
    r.Post("/api/students", handlers.CreateStudent)
    // ... 8 more

    // Program routes (4)
    r.Get("/api/programs", handlers.GetPrograms)
    // ... 3 more

    // College routes (4)
    // Stipend routes (7)
    // Report routes (4)

    port := os.Getenv("PORT")  // 8084
    http.ListenAndServe(":"+port, r)
}
```

---

### 2. `database/db.go` - Database Management (89 lines)

**Purpose**: Establish database connection, configure connection pool, auto-migrate schemas

**Global Variable**:

```go
var DB *gorm.DB  // Singleton database instance
```

**Key Function**:

```go
func Connect() error {
    // 1. Read DATABASE_URL from environment
    dsn := os.Getenv("DATABASE_URL")
    // postgresql://rubadmin:rubpassword@localhost:5432/student_service_db

    // 2. Add SSL mode if missing (for cloud databases)
    if !strings.Contains(dsn, "sslmode=") {
        dsn += "?sslmode=disable"  // Local development
    }

    // 3. Configure GORM logger (Warn level to reduce noise)
    gormLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,  // Log slow queries
            LogLevel:                  logger.Warn,  // Only warnings/errors
            IgnoreRecordNotFoundError: true,
            Colorful:                  false,        // No ANSI colors
        },
    )

    // 4. Open PostgreSQL connection
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })

    // 5. Configure connection pooling for production
    sqlDB, _ := DB.DB()
    sqlDB.SetMaxIdleConns(10)            // 10 idle connections
    sqlDB.SetMaxOpenConns(100)           // Max 100 concurrent
    sqlDB.SetConnMaxLifetime(time.Hour)  // Close after 1 hour

    // 6. Auto-migrate all models (creates tables/columns)
    DB.AutoMigrate(
        &models.Student{},
        &models.College{},
        &models.Program{},
        &models.StipendAllocation{},
        &models.StipendHistory{},
        &models.AuditLog{},
    )

    return nil
}
```

**AutoMigrate Features**:

- Creates tables if they don't exist
- Adds new columns to existing tables
- Creates foreign key constraints
- Creates indexes (unique, composite)
- **Does NOT**: Drop columns, modify types, delete data
- Safe for production (non-destructive)

---

### 3. `handlers/student_handler.go` - Student Operations (217 lines)

**Purpose**: HTTP handlers for student CRUD operations

**Functions** (10 total):

#### `GetStudents()` - List All

```go
func GetStudents(w http.ResponseWriter, r *http.Request) {
    var students []models.Student
    database.DB.Find(&students)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(students)
}
```

#### `CreateStudent()` - Create New

```go
func CreateStudent(w http.ResponseWriter, r *http.Request) {
    var student models.Student
    json.NewDecoder(r.Body).Decode(&student)

    // Validation
    if student.StudentID == "" || student.FirstName == "" ||
       student.LastName == "" || student.Email == "" {
        http.Error(w, "Missing required fields", 400)
        return
    }

    database.DB.Create(&student)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(student)
}
```

#### `SearchStudents()` - Full-Text Search

```go
func SearchStudents(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    searchPattern := "%" + query + "%"

    var students []models.Student
    database.DB.Where(
        "first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR student_id ILIKE ?",
        searchPattern, searchPattern, searchPattern, searchPattern
    ).Find(&students)

    json.NewEncoder(w).Encode(students)
}
```

**Features**:

- Case-insensitive search (ILIKE)
- Searches 4 fields simultaneously
- Partial matching

---

### 4. `handlers/program_handler.go` - Program & College (164 lines)

**Purpose**: Manage academic programs and colleges

**Program Functions** (4):

- `GetPrograms()` - Lists all with College preload
- `CreateProgram()`
- `GetProgramById()`
- `UpdateProgram()`

**College Functions** (4):

- `GetColleges()`
- `CreateCollege()`
- `GetCollegeById()`
- `UpdateCollege()`

**Example - GetPrograms with Relationship**:

```go
func GetPrograms(w http.ResponseWriter, r *http.Request) {
    var programs []models.Program

    // Preload College to avoid N+1 queries
    database.DB.Preload("College").Find(&programs)

    json.NewEncoder(w).Encode(programs)
}
```

**Response**:

```json
[
  {
    "id": 1,
    "code": "BSIT",
    "name": "Bachelor of Science in IT",
    "college": {
      "id": 1,
      "code": "CST",
      "name": "College of Science & Technology"
    },
    "has_stipend": true,
    "stipend_amount": 5000.0
  }
]
```

---

### 5. `handlers/stipend_handler.go` - Stipend Management (268 lines)

**Purpose**: Eligibility checking, allocation, and history tracking

**Key Function - Eligibility Logic**:

```go
func calculateEligibility(student models.Student) models.StipendEligibility {
    eligibility := models.StipendEligibility{StudentID: student.ID}

    // Rule 1: Must be active
    if student.Status != "active" {
        eligibility.IsEligible = false
        eligibility.Reasons = []string{"Student is not active"}
        return eligibility
    }

    // Rule 2: No academic probation
    if student.AcademicStanding == "probation" ||
       student.AcademicStanding == "suspended" {
        eligibility.IsEligible = false
        eligibility.Reasons = []string{"Poor academic standing"}
        return eligibility
    }

    // Rule 3: Minimum GPA 2.0
    if student.GPA < 2.0 {
        eligibility.IsEligible = false
        eligibility.Reasons = []string{"GPA below minimum requirement"}
        return eligibility
    }

    // Rule 4: Program must offer stipend
    if !student.Program.HasStipend {
        eligibility.IsEligible = false
        eligibility.Reasons = []string{"Program does not offer stipend"}
        return eligibility
    }

    // All checks passed!
    eligibility.IsEligible = true
    eligibility.ExpectedAmount = student.Program.StipendAmount
    eligibility.Reasons = []string{"All eligibility criteria met"}

    return eligibility
}
```

**Business Rules**:

1. Active enrollment status
2. Good academic standing (not probation/suspended)
3. Minimum GPA 2.0
4. Program has stipend enabled

---

### 6. `handlers/report_handler.go` - Analytics (170 lines)

**Purpose**: Generate statistical reports and dashboards

**Report Types**:

#### 1. Student Summary

```go
func GenerateStudentSummary(w http.ResponseWriter, r *http.Request) {
    summary := StudentSummary{
        ByCollege: make(map[string]int),
        ByProgram: make(map[string]int),
        ByStatus:  make(map[string]int),
    }

    var students []models.Student
    database.DB.Preload("College").Preload("Program").Find(&students)

    for _, student := range students {
        // Count by status
        switch student.Status {
        case "active":   summary.ActiveStudents++
        case "inactive": summary.InactiveStudents++
        case "graduated": summary.GraduatedStudents++
        }

        // Aggregate by college
        summary.ByCollege[student.College.Name]++

        // Aggregate by program
        summary.ByProgram[student.Program.Name]++

        // Check eligibility
        if calculateEligibility(student).IsEligible {
            summary.StipendEligible++
        }
    }

    json.NewEncoder(w).Encode(summary)
}
```

#### 2. Stipend Statistics

```go
func GenerateStipendStatistics(w http.ResponseWriter, r *http.Request) {
    stats := StipendStatistics{
        ByProgram: make(map[string]float64),
        ByCollege: make(map[string]float64),
    }

    var allocations []models.StipendAllocation
    database.DB.Find(&allocations)

    for _, allocation := range allocations {
        stats.TotalAmount += allocation.Amount

        if allocation.Status == "disbursed" {
            stats.DisbursedAmount += allocation.Amount
        }

        // Fetch student data for aggregation
        var student models.Student
        database.DB.Preload("Program").Preload("College").
                    First(&student, allocation.StudentID)

        stats.ByProgram[student.Program.Name] += allocation.Amount
        stats.ByCollege[student.College.Name] += allocation.Amount
    }

    json.NewEncoder(w).Encode(stats)
}
```

---

### 7. `grpc/server/student_server.go` - gRPC Server (328 lines)

**Purpose**: Implement 10 gRPC methods for internal service communication

**Server Struct**:

```go
type StudentServer struct {
    pb.UnimplementedStudentServiceServer  // Forward compatibility
}
```

**Key Pattern - Proto Conversion**:

```go
func convertStudentToProto(student *models.Student) *pb.StudentResponse {
    return &pb.StudentResponse{
        Id:               uint32(student.ID),
        FirstName:        student.FirstName,
        LastName:         student.LastName,
        StudentId:        student.StudentID,
        Email:            student.Email,
        PhoneNumber:      student.PhoneNumber,
        Gpa:              student.GPA,
        EnrollmentStatus: student.Status,
        // ... 15+ more fields
        CreatedAt:        timestamppb.New(student.CreatedAt),
        UpdatedAt:        timestamppb.New(student.UpdatedAt),
    }
}
```

**Example RPC - GetStudent**:

```go
func (s *StudentServer) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.StudentResponse, error) {
    var student models.Student

    if err := database.DB.Preload("Program").Preload("College").
                         First(&student, req.Id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("student with ID %d not found", req.Id)
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return convertStudentToProto(&student), nil
}
```

**Benefits of gRPC**:

- Type-safe API contracts
- 10x faster than JSON REST
- Binary protocol (HTTP/2)
- Automatic code generation
- Bi-directional streaming support

---

### 8. `grpc/client/user_client.go` - User Service Client (105 lines)

**Purpose**: Communicate with User Service for authentication

**Client Struct**:

```go
type UserGRPCClient struct {
    conn   *grpc.ClientConn
    client pb.UserServiceClient
}
```

**Connection Setup**:

```go
func NewUserGRPCClient() (*UserGRPCClient, error) {
    // Read service URL from environment
    address := os.Getenv("USER_GRPC_URL")
    if address == "" {
        address = "localhost:50052"  // Development fallback
    }

    // Create gRPC connection with 5-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    conn, err := grpc.DialContext(ctx, address,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),  // Wait for connection
    )

    return &UserGRPCClient{
        conn:   conn,
        client: pb.NewUserServiceClient(conn),
    }, nil
}
```

**Methods**:

```go
// Get user information by ID
func (c *UserGRPCClient) GetUser(ctx context.Context, userID uint32) (*pb.UserResponse, error) {
    req := &pb.GetUserRequest{Id: userID}
    return c.client.GetUser(ctx, req)
}

// Validate JWT token
func (c *UserGRPCClient) ValidateToken(ctx context.Context, token string) (*pb.ValidateTokenResponse, error) {
    req := &pb.ValidateTokenRequest{Token: token}
    return c.client.ValidateToken(ctx, req)
}

// Get user role for authorization
func (c *UserGRPCClient) GetUserRole(ctx context.Context, userID uint32) (*pb.UserRoleResponse, error) {
    req := &pb.GetUserRoleRequest{UserId: userID}
    return c.client.GetUserRole(ctx, req)
}
```

---

### 9. `grpc/client/banking_client.go` - Banking Service Client

**Purpose**: Verify bank accounts and manage payment details

**Methods**:

```go
// Get student's bank details
func (c *BankingGRPCClient) GetStudentBankDetails(ctx context.Context, studentID uint32) (*pb.BankDetailsResponse, error)

// Create or update bank details
func (c *BankingGRPCClient) UpsertBankDetails(ctx context.Context,
    studentID, bankID uint32, accountNumber, holderName string) (*pb.BankDetailsResponse, error)

// Verify if bank account is valid
func (c *BankingGRPCClient) VerifyBankAccount(ctx context.Context,
    accountNumber string, bankID uint32) (*pb.VerifyBankAccountResponse, error)
```

---

### 10. `middleware/auth.go` - Authentication Helpers

**Purpose**: JWT validation and user context management

```go
// Extract user ID from request context
func GetUserIDFromContext(ctx context.Context) string {
    if userID, ok := ctx.Value("userID").(string); ok {
        return userID
    }
    return ""
}

// Inject user ID into context
func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, "userID", userID)
}
```

---

### 11. `utils/audit.go` - Audit Logging (55 lines)

**Purpose**: Track all data changes for compliance

```go
func LogAudit(r *http.Request, entityType string, entityID uint,
              action string, changes interface{}) error {

    userID := middleware.GetUserIDFromContext(r.Context())
    changesJSON, _ := json.Marshal(changes)

    log := models.AuditLog{
        EntityType: entityType,  // "student", "stipend", "program"
        EntityID:   entityID,
        Action:     action,      // "create", "update", "delete"
        UserID:     parseUserID(userID),
        Changes:    string(changesJSON),
        IPAddress:  getIPAddress(r),
        UserAgent:  r.UserAgent(),
        Timestamp:  time.Now(),
    }

    return database.DB.Create(&log).Error
}
```

**Usage in Handlers**:

```go
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
    // ... update logic ...

    utils.LogAudit(r, "student", student.ID, "update", map[string]interface{}{
        "old_gpa": oldGPA,
        "new_gpa": newGPA,
        "old_status": oldStatus,
        "new_status": newStatus,
    })
}
```

---

### 12. `proto/student/student.proto` - API Contract

**Purpose**: Define gRPC service contract

```protobuf
syntax = "proto3";

package student;
option go_package = "student_management_service/pb/student";

import "google/protobuf/timestamp.proto";

service StudentService {
  rpc GetStudent(GetStudentRequest) returns (StudentResponse);
  rpc CreateStudent(CreateStudentRequest) returns (StudentResponse);
  rpc UpdateStudent(UpdateStudentRequest) returns (StudentResponse);
  rpc DeleteStudent(DeleteStudentRequest) returns (DeleteResponse);
  rpc ListStudents(ListStudentsRequest) returns (ListStudentsResponse);
  rpc SearchStudents(SearchStudentsRequest) returns (ListStudentsResponse);
  rpc GetStudentsByProgram(GetByProgramRequest) returns (ListStudentsResponse);
  rpc GetStudentsByCollege(GetByCollegeRequest) returns (ListStudentsResponse);
  rpc GetStudentByStudentId(GetStudentByStudentIdRequest) returns (StudentResponse);
  rpc CheckStipendEligibility(StipendEligibilityRequest) returns (StipendEligibilityResponse);
}

message StudentResponse {
  uint32 id = 1;
  string first_name = 2;
  string last_name = 3;
  string student_id = 4;
  string cid = 5;
  string email = 6;
  string phone_number = 7;
  string date_of_birth = 8;
  string gender = 9;
  uint32 program_id = 10;
  uint32 college_id = 11;
  uint32 user_id = 12;
  string permanent_address = 13;
  string current_address = 14;
  string guardian_name = 15;
  string guardian_phone = 16;
  string admission_date = 17;
  string enrollment_status = 18;
  double gpa = 19;
  string academic_standing = 20;
  google.protobuf.Timestamp created_at = 21;
  google.protobuf.Timestamp updated_at = 22;
}

message GetStudentRequest {
  uint32 id = 1;
}

message CreateStudentRequest {
  string first_name = 1;
  string last_name = 2;
  string student_id = 3;
  string cid = 4;
  string email = 5;
  string phone_number = 6;
  string date_of_birth = 7;
  string gender = 8;
  uint32 program_id = 9;
  uint32 college_id = 10;
  uint32 user_id = 11;
  string permanent_address = 12;
  string current_address = 13;
  string guardian_name = 14;
  string guardian_phone = 15;
  string admission_date = 16;
  string enrollment_status = 17;
}

// ... more message definitions
```

**Code Generation**:

```bash
protoc --go_out=. --go-grpc_out=. proto/student/student.proto
```

**Generates**:

- `pb/student/student.pb.go` (message serialization)
- `pb/student/student_grpc.pb.go` (service stubs)

---

## Technologies Used

### Backend Framework

- **Go 1.21+**: Compiled language, high performance
- **Chi Router v5.2.3**: Lightweight, composable HTTP router
  - Middleware support
  - URL parameter extraction
  - Sub-routing capabilities

### Database

- **PostgreSQL 15**: ACID-compliant relational database
- **GORM v1.31.0**: Feature-rich ORM
  - AutoMigrate (schema management)
  - Associations (Has One, Has Many, Belongs To)
  - Soft Deletes
  - Connection pooling
  - Transaction support
  - Hooks (Before/After create/update/delete)

### gRPC & Protocol Buffers

- **grpc-go v1.77.0**: Official Go gRPC implementation
- **protobuf v1.34.1**: Protocol Buffers v3
- **protoc-gen-go v1.31+**: Go protobuf generator
- **protoc-gen-go-grpc v1.3+**: gRPC service generator

### Environment & Deployment

- **godotenv v1.5.1**: Load `.env` files
- **Docker**: Container runtime
- **Docker Compose**: Multi-service orchestration

### Development Tools

- **Protocol Buffer Compiler (protoc) v3.21+**
- **grpcurl**: gRPC command-line client
- **Postman**: HTTP/REST API testing

---

## Testing with Postman

### 1. Create College

```http
POST http://localhost:8084/api/colleges
Content-Type: application/json

{
  "code": "CST",
  "name": "College of Science & Technology",
  "location": "Phuentsholing",
  "is_active": true
}
```

### 2. Create Program

```http
POST http://localhost:8084/api/programs
Content-Type: application/json

{
  "code": "BSIT",
  "name": "Bachelor of Science in Information Technology",
  "college_id": 1,
  "level": "undergraduate",
  "duration_years": 4,
  "duration_semesters": 8,
  "has_stipend": true,
  "stipend_amount": 5000.00,
  "stipend_type": "monthly"
}
```

### 3. Create Student

```http
POST http://localhost:8084/api/students
Content-Type: application/json

{
  "user_id": 1,
  "student_id": "02220123",
  "first_name": "Tshering",
  "last_name": "Dorji",
  "email": "tshering.dorji@rub.edu.bt",
  "phone_number": "+97517123456",
  "cid": "11234567890",
  "program_id": 1,
  "college_id": 1,
  "gpa": 3.5,
  "status": "active",
  "academic_standing": "good"
}
```

### 4. Check Eligibility

```http
GET http://localhost:8084/api/stipend/eligibility/1
```

### 5. Get Reports

```http
GET http://localhost:8084/api/reports/students/summary
GET http://localhost:8084/api/reports/stipend/statistics
```

---

## Summary

### Service Statistics

- **HTTP Endpoints**: 29 REST APIs
- **gRPC Methods**: 10 RPC procedures
- **Database Tables**: 6 entities
- **Total Lines of Code**: ~1,500+ lines
- **Handler Functions**: 25 functions
- **Models**: 6 structs

### Key Features Implemented

✅ Complete student lifecycle management  
✅ Academic program and college administration  
✅ Stipend eligibility checking with business rules  
✅ Allocation and payment tracking  
✅ Comprehensive reporting and analytics  
✅ Hybrid HTTP REST + gRPC architecture  
✅ Service-to-service communication  
✅ Database connection pooling  
✅ Soft delete support  
✅ Audit logging  
✅ Full-text search  
✅ Relationship preloading (N+1 prevention)  
✅ Docker containerization

### Design Patterns Used

- **Repository Pattern**: Database layer abstraction
- **Service Layer**: Business logic in handlers
- **Client-Server**: gRPC communication
- **Dependency Injection**: Database instance
- **Singleton**: Global DB connection

### Security Features (Planned)

- JWT authentication middleware
- Role-based access control
- Audit trail for compliance
- Input validation and sanitization
- SQL injection prevention (GORM parameterization)

### Performance Optimizations

- Connection pooling (10 idle, 100 max)
- gRPC binary protocol (10x faster than JSON)
- Eager loading with Preload()
- Database indexes on foreign keys
- GORM query optimization

---

**Documentation Version**: 1.0  
**Last Updated**: November 26, 2025  
**Service Version**: 1.0.0  
**Author**: Student Management Service Team
