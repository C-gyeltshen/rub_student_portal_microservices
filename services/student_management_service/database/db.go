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
    
    // Check if this is an old schema and needs migration
    var columnExists bool
    DB.Raw(`
        SELECT EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_name='students' AND column_name='name'
        )
    `).Scan(&columnExists)
    
    if columnExists {
        log.Println("Detected old schema. Migrating from old schema to new...")
        
        // Backup old data by creating temporary columns
        DB.Exec(`
            ALTER TABLE students ADD COLUMN IF NOT EXISTS temp_name TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS temp_rub_id TEXT;
            UPDATE students SET temp_name = name, temp_rub_id = rub_id_card_number WHERE temp_name IS NULL;
        `)
        
        // Drop old constraints and columns
        DB.Exec(`ALTER TABLE students DROP COLUMN IF EXISTS name CASCADE;`)
        DB.Exec(`ALTER TABLE students DROP COLUMN IF EXISTS rub_id_card_number CASCADE;`)
        
        // Add new columns without NOT NULL first
        DB.Exec(`
            ALTER TABLE students ADD COLUMN IF NOT EXISTS student_id TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS first_name TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS last_name TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS cid TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS gender TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS permanent_address TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS current_address TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS year_of_study INTEGER DEFAULT 0;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS semester INTEGER DEFAULT 0;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS enrollment_date TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS graduation_date TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'active';
            ALTER TABLE students ADD COLUMN IF NOT EXISTS academic_standing TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS gpa DOUBLE PRECISION DEFAULT 0.0;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS guardian_name TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS guardian_phone_number TEXT;
            ALTER TABLE students ADD COLUMN IF NOT EXISTS guardian_relation TEXT;
        `)
        
        // Migrate data from temp columns
        DB.Exec(`
            UPDATE students 
            SET 
                student_id = COALESCE(temp_rub_id, CONCAT('STU', LPAD(id::text, 6, '0'))),
                first_name = COALESCE(SPLIT_PART(temp_name, ' ', 1), 'Unknown'),
                last_name = COALESCE(NULLIF(SUBSTRING(temp_name FROM POSITION(' ' IN temp_name) + 1), ''), 'Student'),
                cid = temp_rub_id
            WHERE student_id IS NULL OR first_name IS NULL;
        `)
        
        // Ensure email is not null
        DB.Exec(`
            UPDATE students 
            SET email = CONCAT('student', id, '@rub.edu.bt')
            WHERE email IS NULL OR email = '';
        `)
        
        // Drop temp columns
        DB.Exec(`
            ALTER TABLE students DROP COLUMN IF EXISTS temp_name;
            ALTER TABLE students DROP COLUMN IF EXISTS temp_rub_id;
        `)
        
        log.Println("Old schema migration completed.")
    } else {
        // For fresh tables, just ensure student_id column exists
        DB.Exec(`
            DO $$ 
            BEGIN
                IF NOT EXISTS (
                    SELECT 1 FROM information_schema.columns 
                    WHERE table_name='students' AND column_name='student_id'
                ) THEN
                    ALTER TABLE students ADD COLUMN student_id TEXT;
                END IF;
            END $$;
        `)
        
        // Populate any null student_id values
        DB.Exec(`
            UPDATE students 
            SET student_id = CONCAT('STU', LPAD(id::text, 6, '0')) 
            WHERE student_id IS NULL OR student_id = '';
        `)
        
        // Ensure email and names are not null for existing records
        DB.Exec(`
            UPDATE students 
            SET email = CONCAT('student', id, '@rub.edu.bt')
            WHERE email IS NULL OR email = '';
        `)
        
        DB.Exec(`
            UPDATE students 
            SET first_name = 'Unknown'
            WHERE first_name IS NULL OR first_name = '';
        `)
        
        DB.Exec(`
            UPDATE students 
            SET last_name = 'Student'
            WHERE last_name IS NULL OR last_name = '';
        `)
    }
    
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