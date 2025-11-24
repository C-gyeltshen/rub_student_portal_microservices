package database

import (
	"log"
	"os"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"finance_service/models"
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
	gormLogger := logger.Default.LogMode(logger.Warn)

	gormLogger = logger.New(
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

// Migrate runs database migrations for the finance service models
func Migrate() error {
	if err := DB.AutoMigrate(
		&models.Stipend{},
		&models.Deduction{},
		&models.DeductionRule{},
	); err != nil {
		log.Printf("Error running migrations: %v", err)
		return err
	}

	log.Println("Database migrations completed successfully!")
	return nil
}
