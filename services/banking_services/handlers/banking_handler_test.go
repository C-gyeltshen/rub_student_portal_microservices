package handlers

import (
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

func TestGetBanks_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "Bank A", nil, nil, nil).
		AddRow(2, "Bank B", nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/banks", nil)
	w := httptest.NewRecorder()

	GetBanks(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []models.Bank
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Bank A", response[0].Name)
}

func TestGetBanks_DatabaseError(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WillReturnError(gorm.ErrInvalidDB)

	req := httptest.NewRequest("GET", "/banks", nil)
	w := httptest.NewRecorder()

	GetBanks(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateBank_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	bank := models.Bank{Name: "New Bank"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "banks"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "New Bank").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	body, _ := json.Marshal(bank)
	req := httptest.NewRequest("POST", "/banks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateBank(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateBank_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/banks", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	CreateBank(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBankById_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "Test Bank", nil, nil, nil)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	router := chi.NewRouter()
	router.Get("/banks/{id}", GetBankById)

	req := httptest.NewRequest("GET", "/banks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBankById_NotFound(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WithArgs("999", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router := chi.NewRouter()
	router.Get("/banks/{id}", GetBankById)

	req := httptest.NewRequest("GET", "/banks/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateBank_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "Old Bank", nil, nil, nil)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "banks"`)).
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

func TestDeleteBank_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "Test Bank", nil, nil, nil)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "banks" SET "deleted_at"=`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := chi.NewRouter()
	router.Delete("/banks/{id}", DeleteBank)

	req := httptest.NewRequest("DELETE", "/banks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetStudentBankDetails_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	detailsRows := sqlmock.NewRows([]string{"id", "student_id", "bank_id", "account_number", "account_holder_name", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, 123, 1, "1234567890", "John Doe", nil, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "student_bank_details"`)).
		WillReturnRows(detailsRows)

	bankRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Bank")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WillReturnRows(bankRows)

	req := httptest.NewRequest("GET", "/banks/get/student-bank-details", nil)
	w := httptest.NewRecorder()

	GetStudentBankDetails(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateStudentBankDetails_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	details := models.StudentBankDetails{
		StudentID:         123,
		BankID:            1,
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
	}

	bankRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Bank")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WithArgs(uint(1), 1).
		WillReturnRows(bankRows)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "student_bank_details"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	detailsRows := sqlmock.NewRows([]string{"id", "student_id", "bank_id", "account_number", "account_holder_name"}).
		AddRow(1, 123, 1, "1234567890", "John Doe")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "student_bank_details"`)).
		WillReturnRows(detailsRows)

	bankRows2 := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Bank")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WillReturnRows(bankRows2)

	body, _ := json.Marshal(details)
	req := httptest.NewRequest("POST", "/banks/create/student-bank-details", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateStudentBankDetails(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateStudentBankDetails_BankNotFound(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	details := models.StudentBankDetails{
		StudentID:         123,
		BankID:            999,
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "banks"`)).
		WithArgs(uint(999), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	body, _ := json.Marshal(details)
	req := httptest.NewRequest("POST", "/banks/create/student-bank-details", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateStudentBankDetails(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
