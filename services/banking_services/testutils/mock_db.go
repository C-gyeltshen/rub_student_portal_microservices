package testutils

import (
"database/sql"
"banking_services/database"
"github.com/DATA-DOG/go-sqlmock"
"gorm.io/driver/postgres"
"gorm.io/gorm"
)

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

func CleanupMockDB(sqlDB *sql.DB) {
	sqlDB.Close()
}
