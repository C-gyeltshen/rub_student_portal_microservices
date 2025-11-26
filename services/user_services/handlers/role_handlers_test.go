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

func TestGetRoles_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "role_name"}).
		AddRow(1, "Student").
		AddRow(2, "Admin")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/roles", nil)
	w := httptest.NewRecorder()

	GetRoles(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRoles_DatabaseError(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WillReturnError(gorm.ErrInvalidDB)

	req := httptest.NewRequest("GET", "/roles", nil)
	w := httptest.NewRecorder()

	GetRoles(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateRole_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	role := models.UserRole{Name: "New Role"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_roles"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	body, _ := json.Marshal(role)
	req := httptest.NewRequest("POST", "/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateRole(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetRoleById_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "role_name"}).
		AddRow(1, "Student")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	router := chi.NewRouter()
	router.Get("/roles/{id}", GetRoleById)

	req := httptest.NewRequest("GET", "/roles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateRole_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "role_name"}).
		AddRow(1, "Old Role")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_roles"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updatedRole := models.UserRole{Name: "Updated Role"}
	body, _ := json.Marshal(updatedRole)

	router := chi.NewRouter()
	router.Patch("/roles/{id}", UpdateRole)

	req := httptest.NewRequest("PATCH", "/roles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteRole_Success(t *testing.T) {
	db, mock, err := testutils.SetupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "role_name"}).
		AddRow(1, "Test Role")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_roles"`)).
		WithArgs("1", 1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_roles" SET "deleted_at"=`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := chi.NewRouter()
	router.Delete("/roles/{id}", DeleteRole)

	req := httptest.NewRequest("DELETE", "/roles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
