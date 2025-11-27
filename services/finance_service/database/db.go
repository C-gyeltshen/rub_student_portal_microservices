package database

import (
	"finance_service/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect reads the DSN from the environment and establishes the database connection.
func Connect() error {
	var err error

	// 1. Read the Database URL from the environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
	}

	log.Println("Attempting to connect to database using DSN from environment...")

	// 2. Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// 3. Open the database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return err
	}

	// 4. Configure connection pooling for better performance
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully!")
	return nil
}

// CreateExtensions ensures required PostgreSQL extensions are available
func CreateExtensions() error {
	// Create pgcrypto extension if it doesn't exist
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		log.Printf("Error creating pgcrypto extension: %v", err)
		return err
	}
	log.Println("pgcrypto extension verified!")

	// Create uuid-ossp extension if it doesn't exist
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Error creating uuid-ossp extension: %v", err)
		return err
	}
	log.Println("uuid-ossp extension verified!")

	return nil
}

// Migrate runs database migrations for the finance service models
func Migrate() error {
	// First, ensure extensions exist
	if err := CreateExtensions(); err != nil {
		log.Printf("Error creating extensions: %v", err)
		return err
	}

	if err := DB.AutoMigrate(
		&models.Stipend{},
		&models.Deduction{},
		&models.DeductionRule{},
		&models.Transaction{},
		&models.AuditLog{},
	); err != nil {
		log.Printf("Error running migrations: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully!")
	return nil
}
