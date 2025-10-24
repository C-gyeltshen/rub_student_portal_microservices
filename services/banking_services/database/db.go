package database

import (
	"log"
	"os"
	"time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"banking_services/models"
)

var DB *gorm.DB

// Connect reads the DSN from the environment and establishes the database connection.
func Connect() error {
    var err error
    
    // 1. Read the Database URL from the environment variable
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
        // os.Getenv returns an empty string if the variable is not found
    }

    log.Println("Attempting to connect to database using DSN from environment...")

    // 2. Configure GORM logger
    // Set log level to Warn to reduce noise from slow migration queries
    // You can change to logger.Info if you want to see all queries during development
    gormLogger := logger.Default.LogMode(logger.Warn)

    // Optionally adjust the slow query threshold (default is 200ms)
    gormLogger = logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,   // Queries slower than 1 second are logged
            LogLevel:                  logger.Warn,   // Log level: Silent, Error, Warn, Info
            IgnoreRecordNotFoundError: true,          // Don't log ErrRecordNotFound errors
            Colorful:                  false,         // Disable color in Docker logs
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
    sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
    sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

    // 5. AutoMigrate the models
    // GORM will create or update the tables based on your structs
    err = DB.AutoMigrate(&models.Bank{}, &models.StudentBankDetails{})
    if err != nil {
        log.Printf("Error running AutoMigrate: %v", err)
        return err
    }

    log.Println("Database connected and User/Role models migrated successfully.")
    return nil
}