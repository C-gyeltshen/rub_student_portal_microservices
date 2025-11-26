package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"user_services/models"
	"user_services/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUsers_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "first_name", "second_name", "email", "user_role_id"}).
		AddRow(1, "John", "Doe", "john@example.com", 1).
		AddRow(2, "Jane", "Smith", "jane@example.com", 2)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data"`)).
		WillReturnRows(rows)

	roleRows := sqlmock.NewRows([]string{"id", "role_name"}).AddRow(1, "Student").AddRow(2, "Admin")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(roleRows)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	GetUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUsers_DatabaseError(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data"`)).
		WillReturnError(gorm.ErrInvalidDB)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	GetUsers(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateUsers_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	user := models.UserData{
		First_name:  "Test",
		Second_name: "User",
		Email:       "test@example.com",
		UserRoleID:  1,
	}

	roleRows := sqlmock.NewRows([]string{"id", "role_name"}).AddRow(1, "Student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WithArgs(uint(1)).
		WillReturnRows(roleRows)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_data"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateUsers(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateUsers_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/users", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	CreateUsers(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetuserById_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "first_name", "second_name", "email", "user_role_id"}).
		AddRow(1, "John", "Doe", "john@example.com", 1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data"`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	roleRows := sqlmock.NewRows([]string{"id", "role_name"}).AddRow(1, "Student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(roleRows)

	router := chi.NewRouter()
	router.Get("/users/{id}", GetuserById)

	req := httptest.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetuserById_NotFound(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data"`)).
		WithArgs("999", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	router := chi.NewRouter()
	router.Get("/users/{id}", GetuserById)

	req := httptest.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUsersByRoleId_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "first_name", "second_name", "email", "user_role_id"}).
		AddRow(1, "John", "Doe", "john@example.com", 1).
		AddRow(2, "Jane", "Smith", "jane@example.com", 1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_data" WHERE user_role_id = $1`)).
		WithArgs("1").
		WillReturnRows(rows)

	roleRows := sqlmock.NewRows([]string{"id", "role_name"}).AddRow(1, "Student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(roleRows)

	router := chi.NewRouter()
	router.Get("/users/roles/{id}", GetUsersByRoleId)

	req := httptest.NewRequest("GET", "/users/roles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateFinanceOfficer_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	user := models.UserData{
		First_name:  "Finance",
		Second_name: "Officer",
		Email:       "finance@example.com",
	}

	roleRows := sqlmock.NewRows([]string{"id", "role_name"}).AddRow(2, "Finance Officer")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE role_name = $1`)).
		WithArgs("Finance Officer").
		WillReturnRows(roleRows)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_data"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users/finance-officer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateFinanceOfficer(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
