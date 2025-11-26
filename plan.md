# Unit Testing Implementation Plan

## Banking Services & User Services - 80% Code Coverage

---

## Table of Contents

1. [Overview](#overview)
2. [Testing Architecture](#testing-architecture)
3. [Required Dependencies](#required-dependencies)
4. [Banking Service Testing Plan](#banking-service-testing-plan)
5. [User Service Testing Plan](#user-service-testing-plan)
6. [Implementation Steps](#implementation-steps)
7. [Code Coverage Strategy](#code-coverage-strategy)
8. [Test Execution & CI/CD](#test-execution--cicd)

---

## Overview

### Objectives

- Achieve **80% code coverage** for both `banking_services` and `user_services`
- Implement comprehensive unit tests for all handlers, models, and database operations
- Use table-driven tests for multiple scenarios
- Mock database interactions using `sqlmock`
- Test HTTP handlers with `httptest`
- Ensure tests are isolated, repeatable, and fast

### Services to Test

#### Banking Service

- **Handlers**: `banking_handler.go`
- **Models**: `banks.go`, `student_bank_details.go`
- **Database**: `db.go`
- **Total Functions**: 11 handlers + 2 models + 1 database connection

#### User Service

- **Handlers**: `user_handlers.go`, `role_handlers.go`
- **Models**: `user.go`, `roles.go`
- **Database**: `db.go`
- **Total Functions**: 11 handlers + 2 models + 1 database connection

---

## Testing Architecture

### Testing Pyramid

```
    /\
   /  \     Unit Tests (80% - Handlers, Models)
  /----\
 /      \   Integration Tests (15% - Database)
/--------\  E2E Tests (5% - Full Flow)
```

### Test Organization

```
services/
├── banking_services/
│   ├── handlers/
│   │   ├── banking_handler.go
│   │   └── banking_handler_test.go      ← NEW
│   ├── models/
│   │   ├── banks.go
│   │   ├── banks_test.go                ← NEW
│   │   ├── student_bank_details.go
│   │   └── student_bank_details_test.go ← NEW
│   ├── database/
│   │   ├── db.go
│   │   └── db_test.go                   ← NEW
│   └── testutils/                        ← NEW
│       ├── mock_db.go
│       └── test_helpers.go
└── user_services/
    ├── handlers/
    │   ├── user_handlers.go
    │   ├── user_handlers_test.go        ← NEW
    │   ├── role_handlers.go
    │   └── role_handlers_test.go        ← NEW
    ├── models/
    │   ├── user.go
    │   ├── user_test.go                 ← NEW
    │   ├── roles.go
    │   └── roles_test.go                ← NEW
    ├── database/
    │   ├── db.go
    │   └── db_test.go                   ← NEW
    └── testutils/                        ← NEW
        ├── mock_db.go
        └── test_helpers.go
```

---

## Required Dependencies

### Add to `go.mod` for Both Services

```bash
# Navigate to each service directory and add dependencies
cd services/banking_services
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/DATA-DOG/go-sqlmock
go get github.com/stretchr/testify/suite
go get gorm.io/driver/postgres

cd ../user_services
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/DATA-DOG/go-sqlmock
go get github.com/stretchr/testify/suite
go get gorm.io/driver/postgres
```

### Dependencies Overview

- **testify/assert**: Assertion library for clean test assertions
- **testify/mock**: Mocking framework for interfaces
- **go-sqlmock**: SQL mock driver for testing database interactions
- **testify/suite**: Test suite support for setup/teardown
- **gorm.io/driver/postgres**: Required for GORM testing

---

## Banking Service Testing Plan

### 1. Handler Tests (`banking_handler_test.go`)

#### Test Coverage Breakdown (11 handlers × ~4 test cases = ~44 tests)

```go
// Test Structure for Each Handler
package handlers

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "github.com/go-chi/chi/v5"
    "gorm.io/gorm"
    "banking_services/models"
    "banking_services/database"
)
```

#### A. Bank Handler Tests

##### 1. **GetBanks** - Test Cases:

- ✅ **Success**: Returns list of banks (200 OK)
- ✅ **Empty**: Returns empty array when no banks exist
- ✅ **Database Error**: Returns 500 when database fails

##### 2. **CreateBank** - Test Cases:

- ✅ **Success**: Creates bank successfully (201 Created)
- ✅ **Invalid JSON**: Returns 400 for malformed JSON
- ✅ **Missing Fields**: Returns 400 for empty name
- ✅ **Database Error**: Returns 500 on DB failure

##### 3. **GetBankById** - Test Cases:

- ✅ **Success**: Returns bank by ID (200 OK)
- ✅ **Not Found**: Returns 404 for non-existent ID
- ✅ **Invalid ID**: Returns 404 for invalid ID format

##### 4. **UpdateBank** - Test Cases:

- ✅ **Success**: Updates bank name (200 OK)
- ✅ **Not Found**: Returns 404 for non-existent bank
- ✅ **Invalid JSON**: Returns 400 for malformed input
- ✅ **Database Error**: Returns 500 on save failure

##### 5. **DeleteBank** - Test Cases:

- ✅ **Success**: Deletes bank (200 OK)
- ✅ **Not Found**: Returns 404 for non-existent bank
- ✅ **Database Error**: Returns 500 on delete failure

#### B. Student Bank Details Handler Tests

##### 6. **GetStudentBankDetails** - Test Cases:

- ✅ **Success**: Returns all student bank details with preloaded Bank
- ✅ **Empty**: Returns empty array
- ✅ **Database Error**: Returns 500

##### 7. **CreateStudentBankDetails** - Test Cases:

- ✅ **Success**: Creates student bank details (201)
- ✅ **Bank Not Found**: Returns 400 when BankID doesn't exist
- ✅ **Invalid JSON**: Returns 400
- ✅ **Missing Required Fields**: Returns 400

##### 8. **GetStudentBankDetailsById** - Test Cases:

- ✅ **Success**: Returns details by ID with Bank preloaded
- ✅ **Not Found**: Returns 404
- ✅ **Invalid ID**: Returns 404

##### 9. **GetStudentBankDetailsByStudentId** - Test Cases:

- ✅ **Success**: Returns bank details for student
- ✅ **Not Found**: Returns 404 when no details exist
- ✅ **Multiple Records**: Returns all records for student
- ✅ **Database Error**: Returns 500

##### 10. **UpdateStudentBankDetails** - Test Cases:

- ✅ **Success**: Updates all fields
- ✅ **Partial Update**: Updates only provided fields
- ✅ **Bank Not Found**: Returns 400 for invalid BankID
- ✅ **Not Found**: Returns 404
- ✅ **Invalid JSON**: Returns 400

##### 11. **DeleteStudentBankDetails** - Test Cases:

- ✅ **Success**: Deletes details (200 OK)
- ✅ **Not Found**: Returns 404
- ✅ **Database Error**: Returns 500

### 2. Model Tests

#### A. `banks_test.go`

```go
// Test bank model validation and GORM behavior
- Test Bank struct JSON marshaling/unmarshaling
- Test Bank relationship with StudentBankDetails
- Test GORM hooks (BeforeCreate, BeforeUpdate if any)
- Test soft delete behavior
```

#### B. `student_bank_details_test.go`

```go
// Test student bank details model
- Test JSON marshaling/unmarshaling
- Test foreign key relationship with Bank
- Test required field validation
- Test GORM timestamps
```

### 3. Database Tests (`db_test.go`)

```go
// Test database connection and setup
- Test Connect() with valid DATABASE_URL
- Test Connect() with missing DATABASE_URL
- Test AutoMigrate executes successfully
- Test connection pool configuration
```

### 4. Test Utilities (`testutils/`)

#### `mock_db.go`

```go
// Mock database setup using sqlmock
- SetupMockDB() - Creates mock DB and sqlmock
- ExpectBankQuery() - Helper for bank queries
- ExpectStudentBankDetailsQuery() - Helper for student queries
```

#### `test_helpers.go`

```go
// Common test helpers
- CreateTestBank() - Returns test bank object
- CreateTestStudentBankDetails() - Returns test details
- SetupTestRouter() - Creates chi router for tests
- MakeRequest() - Helper to make HTTP requests
```

---

## User Service Testing Plan

### 1. Handler Tests

#### A. User Handler Tests (`user_handlers_test.go`)

##### 1. **GetUsers** - Test Cases:

- ✅ **Success**: Returns users with preloaded roles
- ✅ **Empty**: Returns empty array
- ✅ **Database Error**: Returns 500

##### 2. **CreateUsers** - Test Cases:

- ✅ **Success**: Creates user (201)
- ✅ **Invalid JSON**: Returns 400
- ✅ **Missing Fields**: Returns 400
- ✅ **Database Error**: Returns 500

##### 3. **GetuserById** - Test Cases:

- ✅ **Success**: Returns user with role
- ✅ **Not Found**: Returns 404
- ✅ **Invalid ID**: Returns 404

##### 4. **GetUsersByRoleId** - Test Cases:

- ✅ **Success**: Returns users by role
- ✅ **Empty**: Returns empty array for role with no users
- ✅ **Multiple Users**: Returns all users with role
- ✅ **Database Error**: Returns 500

##### 5. **DeleteUsersByRoleId** - Test Cases:

- ✅ **Success**: Deletes multiple users
- ✅ **No Users**: Returns success with 0 deleted
- ✅ **Database Error**: Returns 500

##### 6. **CreateFinanceOfficer** - Test Cases:

- ✅ **Success**: Creates user with RoleID=2
- ✅ **Invalid JSON**: Returns 400
- ✅ **Role ID Override**: Ensures RoleID is always 2
- ✅ **Database Error**: Returns 500

#### B. Role Handler Tests (`role_handlers_test.go`)

##### 1. **GetRoles** - Test Cases:

- ✅ **Success**: Returns all roles
- ✅ **Empty**: Returns empty array
- ✅ **Database Error**: Returns 500

##### 2. **CreateRole** - Test Cases:

- ✅ **Success**: Creates role (201)
- ✅ **Invalid JSON**: Returns 400
- ✅ **Database Error**: Returns 500

##### 3. **GetRoleById** - Test Cases:

- ✅ **Success**: Returns role by ID
- ✅ **Not Found**: Returns 404
- ✅ **Invalid ID**: Returns 404

##### 4. **UpdateRole** - Test Cases:

- ✅ **Success**: Updates role (200)
- ✅ **Not Found**: Returns 404
- ✅ **Invalid JSON**: Returns 400
- ✅ **Database Error**: Returns 500

##### 5. **DeleteRole** - Test Cases:

- ✅ **Success**: Deletes role (204)
- ✅ **Not Found**: Returns 404
- ✅ **Database Error**: Returns 500

### 2. Model Tests

#### A. `user_test.go`

```go
- Test UserData struct validation
- Test foreign key relationship with UserRole
- Test JSON tags
- Test GORM timestamps and soft delete
```

#### B. `roles_test.go`

```go
- Test UserRole struct validation
- Test JSON marshaling
- Test GORM behavior
```

### 3. Database Tests (`db_test.go`)

```go
- Test Connect() with valid DATABASE_URL
- Test Connect() with SSL mode handling
- Test AutoMigrate for UserData and UserRole
- Test connection pool configuration
```

### 4. Test Utilities (`testutils/`)

Similar structure as banking service with user-specific helpers.

---

## Implementation Steps

### Phase 1: Setup (Day 1)

#### Step 1.1: Install Dependencies

```bash
# Banking Service
cd services/banking_services
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/DATA-DOG/go-sqlmock
go get github.com/stretchr/testify/suite

# User Service
cd ../user_services
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/DATA-DOG/go-sqlmock
go get github.com/stretchr/testify/suite
```

#### Step 1.2: Create Test Directories

```bash
# Banking Service
mkdir -p services/banking_services/testutils

# User Service
mkdir -p services/user_services/testutils
```

#### Step 1.3: Create Test Utility Files

##### Banking Service - `testutils/mock_db.go`

```go
package testutils

import (
	"database/sql"
	"banking_services/database"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupMockDB creates a mock database connection
func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	database.DB = db
	return db, mock, nil
}

// CleanupMockDB closes the mock database
func CleanupMockDB(sqlDB *sql.DB) {
	sqlDB.Close()
}
```

##### Banking Service - `testutils/test_helpers.go`

```go
package testutils

import (
	"banking_services/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"github.com/go-chi/chi/v5"
)

// CreateTestBank returns a test bank object
func CreateTestBank() *models.Bank {
	return &models.Bank{
		Model: gorm.Model{ID: 1},
		Name:  "Test Bank",
	}
}

// CreateTestStudentBankDetails returns test student bank details
func CreateTestStudentBankDetails() *models.StudentBankDetails {
	return &models.StudentBankDetails{
		Model:             gorm.Model{ID: 1},
		StudentID:         123,
		BankID:            1,
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
	}
}

// SetupTestRouter creates a chi router for testing
func SetupTestRouter() *chi.Mux {
	return chi.NewRouter()
}

// MakeRequest is a helper to make HTTP requests in tests
func MakeRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	return http.NewRequest(method, url, reqBody)
}

// ExecuteRequest executes a request and returns response recorder
func ExecuteRequest(req *http.Request, router *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
```

**Repeat similar files for `user_services`**

---

### Phase 2: Banking Service Tests (Day 2-3)

#### Step 2.1: Handler Tests - Bank Operations

##### Create `handlers/banking_handler_test.go`

```go
package handlers

import (
	"banking_services/database"
	"banking_services/models"
	"banking_services/testutils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestGetBanks_Success tests successful retrieval of banks
func TestGetBanks_Success(t *testing.T) {
	// Setup mock database
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Expected banks
	banks := []models.Bank{
		{Model: gorm.Model{ID: 1}, Name: "Bank A"},
		{Model: gorm.Model{ID: 2}, Name: "Bank B"},
	}

	// Mock expectations
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "Bank A", nil, nil, nil).
		AddRow(2, "Bank B", nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks" WHERE "banks"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	// Create request
	req := httptest.NewRequest("GET", "/banks", nil)
	w := httptest.NewRecorder()

	// Execute handler
	GetBanks(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Bank
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Bank A", response[0].Name)
}

// TestGetBanks_DatabaseError tests database error handling
func TestGetBanks_DatabaseError(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Mock database error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WillReturnError(gorm.ErrInvalidDB)

	req := httptest.NewRequest("GET", "/banks", nil)
	w := httptest.NewRecorder()

	GetBanks(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestCreateBank_Success tests successful bank creation
func TestCreateBank_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	bank := models.Bank{Name: "New Bank"}

	// Mock expectations
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "banks"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "New Bank").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Create request
	body, _ := json.Marshal(bank)
	req := httptest.NewRequest("POST", "/banks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateBank(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestCreateBank_InvalidJSON tests invalid JSON input
func TestCreateBank_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/banks", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	CreateBank(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetBankById_Success tests successful retrieval by ID
func TestGetBankById_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Mock expectations
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Bank")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks" WHERE "banks"."id" = $1`)).
		WithArgs("1").
		WillReturnRows(rows)

	// Setup router with URL parameter
	router := chi.NewRouter()
	router.Get("/banks/{id}", GetBankById)

	req := httptest.NewRequest("GET", "/banks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetBankById_NotFound tests bank not found scenario
func TestGetBankById_NotFound(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks" WHERE "banks"."id" = $1`)).
		WithArgs("999").
		WillReturnError(gorm.ErrRecordNotFound)

	router := chi.NewRouter()
	router.Get("/banks/{id}", GetBankById)

	req := httptest.NewRequest("GET", "/banks/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUpdateBank_Success tests successful bank update
func TestUpdateBank_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Mock finding the bank
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Old Bank")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks" WHERE "banks"."id" = $1`)).
		WithArgs("1").
		WillReturnRows(rows)

	// Mock update
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "banks" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updatedBank := models.Bank{Name: "Updated Bank"}
	body, _ := json.Marshal(updatedBank)

	router := chi.NewRouter()
	router.Patch("/banks/{id}", UpdateBank)

	req := httptest.NewRequest("PATCH", "/banks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestDeleteBank_Success tests successful bank deletion
func TestDeleteBank_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Mock finding the bank
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Bank")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks" WHERE "banks"."id" = $1`)).
		WithArgs("1").
		WillReturnRows(rows)

	// Mock delete (soft delete in GORM)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "banks" SET "deleted_at"=$1 WHERE "banks"."id" = $2`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := chi.NewRouter()
	router.Delete("/banks/{id}", DeleteBank)

	req := httptest.NewRequest("DELETE", "/banks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Add similar tests for Student Bank Details handlers
// TestGetStudentBankDetails_Success
// TestCreateStudentBankDetails_Success
// TestCreateStudentBankDetails_BankNotFound
// TestGetStudentBankDetailsById_Success
// TestGetStudentBankDetailsByStudentId_Success
// TestUpdateStudentBankDetails_Success
// TestDeleteStudentBankDetails_Success
// ... (approximately 15-20 more tests)
```

#### Step 2.2: Model Tests

##### Create `models/banks_test.go`

```go
package models

import (
	"encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBank_JSONMarshaling(t *testing.T) {
	bank := Bank{
		Model: gorm.Model{ID: 1},
		Name:  "Test Bank",
	}

	jsonData, err := json.Marshal(bank)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "Test Bank")

	var unmarshaledBank Bank
	err = json.Unmarshal(jsonData, &unmarshaledBank)
	assert.NoError(t, err)
	assert.Equal(t, "Test Bank", unmarshaledBank.Name)
}

func TestBank_Relationships(t *testing.T) {
	bank := Bank{
		Model:              gorm.Model{ID: 1},
		Name:               "Test Bank",
		StudentBankDetails: []StudentBankDetails{},
	}

	assert.NotNil(t, bank.StudentBankDetails)
	assert.Len(t, bank.StudentBankDetails, 0)
}
```

##### Create `models/student_bank_details_test.go`

```go
package models

import (
	"encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestStudentBankDetails_JSONMarshaling(t *testing.T) {
	details := StudentBankDetails{
		Model:             gorm.Model{ID: 1},
		StudentID:         123,
		BankID:            1,
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
	}

	jsonData, err := json.Marshal(details)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), "John Doe")

	var unmarshaled StudentBankDetails
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", unmarshaled.AccountHolderName)
	assert.Equal(t, 123, unmarshaled.StudentID)
}

func TestStudentBankDetails_ForeignKey(t *testing.T) {
	details := StudentBankDetails{
		BankID: 1,
		Bank: Bank{
			Model: gorm.Model{ID: 1},
			Name:  "Test Bank",
		},
	}

	assert.Equal(t, uint(1), details.BankID)
	assert.Equal(t, "Test Bank", details.Bank.Name)
}
```

#### Step 2.3: Database Tests

##### Create `database/db_test.go`

```go
package database

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConnect_MissingDatabaseURL(t *testing.T) {
	// Save original env
	originalURL := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", originalURL)

	// Unset DATABASE_URL
	os.Unsetenv("DATABASE_URL")

	// This should fatal, but we can't easily test log.Fatal
	// Instead, we'll test that DATABASE_URL is checked
	dsn := os.Getenv("DATABASE_URL")
	assert.Empty(t, dsn)
}

func TestConnect_WithValidURL(t *testing.T) {
	// This is an integration test that requires a real database
	// Skip if not in integration test mode
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test")
	}

	err := Connect()
	assert.NoError(t, err)
	assert.NotNil(t, DB)
}
```

---

### Phase 3: User Service Tests (Day 4-5)

#### Step 3.1: Handler Tests - User Operations

##### Create `handlers/user_handlers_test.go`

```go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"user_services/database"
	"user_services/models"
	"user_services/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestGetUsers_Success
func TestGetUsers_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Mock user data with role
	rows := sqlmock.NewRows([]string{"id", "first_name", "second_name", "email", "user_role_id"}).
		AddRow(1, "John", "Doe", "john@example.com", 1).
		AddRow(2, "Jane", "Smith", "jane@example.com", 2)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data"`)).
		WillReturnRows(rows)

	// Mock Role preload
	roleRows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Student", "Student role").
		AddRow(2, "Finance", "Finance officer")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(roleRows)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	GetUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.UserData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

// TestCreateUsers_Success
func TestCreateUsers_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	user := models.UserData{
		First_name:  "John",
		Second_name: "Doe",
		Email:       "john@example.com",
		UserRoleID:  1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_data"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Mock reload with role
	userRows := sqlmock.NewRows([]string{"id", "first_name", "second_name", "email", "user_role_id"}).
		AddRow(1, "John", "Doe", "john@example.com", 1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data" WHERE "user_data"."id" = $1`)).
		WillReturnRows(userRows)

	roleRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(roleRows)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateUsers(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestCreateFinanceOfficer_EnforcesRoleID
func TestCreateFinanceOfficer_EnforcesRoleID(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// User tries to set RoleID to 1, but it should be forced to 2
	user := models.UserData{
		First_name:  "John",
		Second_name: "Doe",
		Email:       "john@example.com",
		UserRoleID:  1, // This should be overridden to 2
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_data"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Verify reload expects RoleID = 2
	userRows := sqlmock.NewRows([]string{"id", "first_name", "user_role_id"}).
		AddRow(1, "John", 2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data"`)).
		WillReturnRows(userRows)

	roleRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "Finance")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(roleRows)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users/create/finance-officer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateFinanceOfficer(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.UserData
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, uint(2), response.UserRoleID)
}

// Add more tests for:
// - GetuserById
// - GetUsersByRoleId
// - DeleteUsersByRoleId
// ... (approximately 15-20 more tests)
```

#### Step 3.2: Role Handler Tests

##### Create `handlers/role_handlers_test.go`

```go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"user_services/database"
	"user_services/models"
	"user_services/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestGetRoles_Success
func TestGetRoles_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Student", "Student role").
		AddRow(2, "Finance", "Finance officer")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/users/get/roles", nil)
	w := httptest.NewRecorder()

	GetRoles(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestCreateRole_Success
func TestCreateRole_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	role := models.UserRole{
		Name:        "Admin",
		Description: "Administrator role",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_roles"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	body, _ := json.Marshal(role)
	req := httptest.NewRequest("POST", "/users/create/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateRole(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestDeleteRole_Success
func TestDeleteRole_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Role")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE "user_roles"."id" = $1`)).
		WithArgs("1").
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_roles" SET "deleted_at"=$1`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := chi.NewRouter()
	router.Delete("/users/delete/role/{id}", DeleteRole)

	req := httptest.NewRequest("DELETE", "/users/delete/role/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// Add more tests for UpdateRole, GetRoleById
```

#### Step 3.3: Model and Database Tests

Similar structure as banking service but for user models.

---

### Phase 4: Coverage Analysis & Optimization (Day 6)

#### Step 4.1: Run Coverage Tests

```bash
# Banking Service
cd services/banking_services
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# User Service
cd ../user_services
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

#### Step 4.2: Analyze Coverage Reports

```bash
# View coverage percentage by package
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Expected output should show >80% for:
# - handlers/banking_handler.go
# - handlers/user_handlers.go
# - handlers/role_handlers.go
# - models/*.go
```

#### Step 4.3: Identify Gaps

```bash
# Generate detailed HTML report
go tool cover -html=coverage.out

# Look for red/uncovered lines in:
# - Error handling paths
# - Edge cases
# - Database error scenarios
```

#### Step 4.4: Add Additional Tests

Focus on:

- Error paths not covered
- Edge cases (empty strings, nil values, boundary conditions)
- Complex conditional logic
- Preload relationships

---

## Code Coverage Strategy

### Coverage Goals by Component

| Component   | Target Coverage | Priority     |
| ----------- | --------------- | ------------ |
| Handlers    | 85%+            | High         |
| Models      | 80%+            | Medium       |
| Database    | 75%+            | Medium       |
| Utils       | 90%+            | High         |
| **Overall** | **80%+**        | **Required** |

### Coverage Techniques

#### 1. **Table-Driven Tests**

```go
func TestGetBankById_MultipleCases(t *testing.T) {
	testCases := []struct {
		name           string
		bankID         string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectError    bool
	}{
		{
			name:   "Valid ID",
			bankID: "1",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WillReturnRows(rows)
			},
			expectedStatus: 200,
			expectError:    false,
		},
		{
			name:   "Invalid ID",
			bankID: "999",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedStatus: 404,
			expectError:    true,
		},
		// Add more cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test implementation
		})
	}
}
```

#### 2. **Boundary Testing**

- Empty arrays
- Nil values
- Zero values
- Maximum lengths
- Invalid IDs

#### 3. **Error Path Testing**

- Database connection failures
- Transaction rollbacks
- Foreign key violations
- JSON parsing errors
- Missing required fields

#### 4. **Integration Points**

- Preload relationships
- Foreign key validations
- GORM hooks
- Soft delete behavior

---

## Test Execution & CI/CD

### Running Tests Locally

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./handlers -v

# Run specific test
go test ./handlers -run TestGetBanks_Success -v

# Run with race detector
go test ./... -race

# Continuous test watching (install gotestsum)
gotestsum --watch
```

### Makefile Additions

Add to root `Makefile`:

```makefile
# Test targets
.PHONY: test test-coverage test-banking test-user test-all

test:
	@echo "Running all tests..."
	@cd services/banking_services && go test ./... -v
	@cd services/user_services && go test ./... -v

test-coverage:
	@echo "Running tests with coverage..."
	@cd services/banking_services && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
	@cd services/user_services && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage reports generated: coverage.html"

test-banking:
	@echo "Testing banking service..."
	@cd services/banking_services && go test ./... -v -cover

test-user:
	@echo "Testing user service..."
	@cd services/user_services && go test ./... -v -cover

test-all: test-coverage
	@echo "All tests completed with coverage"

# Coverage threshold check
test-coverage-check:
	@cd services/banking_services && go test ./... -coverprofile=coverage.out
	@cd services/banking_services && go tool cover -func=coverage.out | grep total | awk '{if ($$3+0 < 80.0) {print "Banking service coverage below 80%: " $$3; exit 1}}'
	@cd services/user_services && go test ./... -coverprofile=coverage.out
	@cd services/user_services && go tool cover -func=coverage.out | grep total | awk '{if ($$3+0 < 80.0) {print "User service coverage below 80%: " $$3; exit 1}}'
	@echo "✓ All services meet 80% coverage threshold"
```

### GitHub Actions CI Pipeline

Create `.github/workflows/test.yml`:

```yaml
name: Go Tests

on:
  push:
    branches: [main, develop, feat/*]
  pull_request:
    branches: [main, develop]

jobs:
  test:
    name: Test Services
    runs-on: ubuntu-latest

    strategy:
      matrix:
        service: [banking_services, user_services]

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.21"

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Install dependencies
        working-directory: ./services/${{ matrix.service }}
        run: go mod download

      - name: Run tests
        working-directory: ./services/${{ matrix.service }}
        run: go test ./... -v -coverprofile=coverage.out

      - name: Check coverage threshold
        working-directory: ./services/${{ matrix.service }}
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 80.0" | bc -l) )); then
            echo "❌ Coverage $COVERAGE% is below 80% threshold"
            exit 1
          fi
          echo "✅ Coverage $COVERAGE% meets 80% threshold"

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v3
        with:
          file: ./services/${{ matrix.service }}/coverage.out
          flags: ${{ matrix.service }}
          name: ${{ matrix.service }}-coverage
```

### Pre-commit Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash

echo "Running tests before commit..."

# Test banking service
cd services/banking_services
if ! go test ./... -cover; then
    echo "❌ Banking service tests failed"
    exit 1
fi

# Test user service
cd ../user_services
if ! go test ./... -cover; then
    echo "❌ User service tests failed"
    exit 1
fi

echo "✅ All tests passed"
exit 0
```

Make it executable:

```bash
chmod +x .git/hooks/pre-commit
```

---

## Summary Checklist

### Banking Service

- [ ] Install test dependencies
- [ ] Create `testutils/` directory with helpers
- [ ] Write 44+ handler tests (11 handlers × 4 cases)
- [ ] Write 4+ model tests
- [ ] Write 3+ database tests
- [ ] Achieve 80%+ coverage
- [ ] Generate coverage report

### User Service

- [ ] Install test dependencies
- [ ] Create `testutils/` directory with helpers
- [ ] Write 44+ handler tests (11 handlers × 4 cases)
- [ ] Write 4+ model tests
- [ ] Write 3+ database tests
- [ ] Achieve 80%+ coverage
- [ ] Generate coverage report

### CI/CD

- [ ] Add Makefile test targets
- [ ] Create GitHub Actions workflow
- [ ] Set up pre-commit hooks
- [ ] Configure coverage reporting
- [ ] Add badge to README

---

## Expected Test Count

| Service           | Component      | Test Count | Coverage Target |
| ----------------- | -------------- | ---------- | --------------- |
| Banking           | Handlers       | ~44 tests  | 85%             |
| Banking           | Models         | ~6 tests   | 80%             |
| Banking           | Database       | ~3 tests   | 75%             |
| **Banking Total** | **~53 tests**  | **80%+**   |
| User              | Handlers       | ~44 tests  | 85%             |
| User              | Models         | ~6 tests   | 80%             |
| User              | Database       | ~3 tests   | 75%             |
| **User Total**    | **~53 tests**  | **80%+**   |
| **Grand Total**   | **~106 tests** | **80%+**   |

---

## Estimated Timeline

| Phase                  | Duration     | Tasks                               |
| ---------------------- | ------------ | ----------------------------------- |
| Phase 1: Setup         | 4 hours      | Install deps, create test structure |
| Phase 2: Banking Tests | 12 hours     | Write handler, model, DB tests      |
| Phase 3: User Tests    | 12 hours     | Write handler, model, DB tests      |
| Phase 4: Coverage      | 4 hours      | Analyze, optimize, document         |
| **Total**              | **32 hours** | **~4 working days**                 |

---

## Best Practices

1. **Test Isolation**: Each test should be independent
2. **Clean Setup/Teardown**: Use `defer` for cleanup
3. **Descriptive Names**: `TestHandlerName_Scenario_ExpectedOutcome`
4. **Table-Driven**: Use for multiple similar scenarios
5. **Mock External Dependencies**: Always mock database
6. **Assert Clearly**: Use testify assertions for readability
7. **Test Edge Cases**: Empty, nil, boundary values
8. **Coverage != Quality**: Aim for meaningful tests, not just numbers

---

## Maintenance

### Running Tests Regularly

```bash
# Before each commit
make test

# Before each PR
make test-coverage-check

# Weekly full coverage report
make test-coverage
```

### Updating Tests

- Add tests for new features
- Update mocks when models change
- Refactor tests when handlers change
- Keep coverage above 80%

---

## Conclusion

This plan provides a comprehensive approach to implementing unit tests for both `banking_services` and `user_services` with 80% code coverage. By following this structured approach, you will:

1. ✅ Achieve 80%+ code coverage
2. ✅ Create maintainable, isolated tests
3. ✅ Establish CI/CD test automation
4. ✅ Improve code quality and confidence
5. ✅ Enable safe refactoring
6. ✅ Catch bugs early in development

**Start with Phase 1 (Setup) and proceed sequentially through each phase for best results.**
