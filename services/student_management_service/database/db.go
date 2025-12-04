package database

import (
	"log"
	"os"
	"strings"
	"student_management_service/models"
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
        // os.Getenv returns an empty string if the variable is not found
    }

    // 2. Add SSL mode if not present
    // Use sslmode=disable for local/Docker, sslmode=require for cloud
    if !strings.Contains(dsn, "sslmode=") {
        if strings.Contains(dsn, "?") {
            dsn += "&sslmode=disable"
        } else {
            dsn += "?sslmode=disable"
        }
    }

    log.Println("Attempting to connect to database using DSN from environment...")

    // 3. Configure GORM logger
    // Set log level to Warn to reduce noise from slow migration queries
    // You can change to logger.Info if you want to see all queries during development
    gormLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,   // Queries slower than 1 second are logged
            LogLevel:                  logger.Warn,   // Log level: Silent, Error, Warn, Info
            IgnoreRecordNotFoundError: true,          // Don't log ErrRecordNotFound errors
            Colorful:                  false,         // Disable color in Docker logs
        },
    )

    // 4. Open the database connection
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })
    if err != nil {
        log.Printf("Error connecting to database: %v", err)
        return err
    }

    // 5. Configure connection pooling for better performance
    sqlDB, err := DB.DB()
    if err != nil {
        log.Printf("Error getting database instance: %v", err)
        return err
    }

    // Set connection pool settings
    sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
    sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

    // 6. AutoMigrate the models
    // GORM will create or update the tables based on your structs
    // Order matters: Create tables that are referenced by foreign keys first
    log.Println("Running AutoMigrate for colleges, programs, and students...")
    err = DB.AutoMigrate(
        &models.College{},
        &models.Program{},
        &models.Student{},
    )
    if err != nil {
        log.Printf("Error running AutoMigrate: %v", err)
        return err
    }

    // Create stipend tables if they don't exist
    // Drop existing tables first to ensure clean schema (only if they exist)
    DB.Exec(`DROP TABLE IF EXISTS stipend_histories CASCADE;`)
    DB.Exec(`DROP TABLE IF EXISTS stipend_allocations CASCADE;`)
    
    // Create stipend_allocations table
    DB.Exec(`
        CREATE TABLE IF NOT EXISTS stipend_allocations (
            id SERIAL PRIMARY KEY,
            allocation_id VARCHAR(255) UNIQUE NOT NULL,
            student_id INTEGER NOT NULL,
            amount DECIMAL(10,2) NOT NULL,
            allocation_date VARCHAR(255),
            status VARCHAR(50) DEFAULT 'pending',
            approved_by INTEGER,
            approval_date VARCHAR(255),
            semester INTEGER,
            academic_year VARCHAR(255),
            remarks TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP,
            FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
        );
    `)

    // Create stipend_histories table
    DB.Exec(`
        CREATE TABLE IF NOT EXISTS stipend_histories (
            id SERIAL PRIMARY KEY,
            transaction_id VARCHAR(255) UNIQUE NOT NULL,
            student_id INTEGER NOT NULL,
            allocation_id INTEGER,
            amount DECIMAL(10,2) NOT NULL,
            payment_date VARCHAR(255),
            transaction_status VARCHAR(50),
            payment_method VARCHAR(50),
            bank_reference VARCHAR(255),
            remarks TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP,
            FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
            FOREIGN KEY (allocation_id) REFERENCES stipend_allocations(id) ON DELETE SET NULL
        );
    `)

    // Create indexes
    DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_stipend_allocations_student_id ON stipend_allocations(student_id);
        CREATE INDEX IF NOT EXISTS idx_stipend_allocations_status ON stipend_allocations(status);
        CREATE INDEX IF NOT EXISTS idx_stipend_allocations_allocation_id ON stipend_allocations(allocation_id);
        CREATE INDEX IF NOT EXISTS idx_stipend_allocations_deleted_at ON stipend_allocations(deleted_at);
        
        CREATE INDEX IF NOT EXISTS idx_stipend_histories_student_id ON stipend_histories(student_id);
        CREATE INDEX IF NOT EXISTS idx_stipend_histories_allocation_id ON stipend_histories(allocation_id);
        CREATE INDEX IF NOT EXISTS idx_stipend_histories_transaction_id ON stipend_histories(transaction_id);
        CREATE INDEX IF NOT EXISTS idx_stipend_histories_transaction_status ON stipend_histories(transaction_status);
        CREATE INDEX IF NOT EXISTS idx_stipend_histories_deleted_at ON stipend_histories(deleted_at);
    `)

    log.Println("Database connected and all models migrated successfully.")
    return nil
}