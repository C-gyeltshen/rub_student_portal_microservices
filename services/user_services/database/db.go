package database

import (
	"log"
	"os" 
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"user_services/models" 
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

    // 2. Open the database connection
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("Error connecting to database: %v", err)
        return err
    }

    // 3. AutoMigrate the models
    // GORM will create or update the tables based on your structs
    err = DB.AutoMigrate(&models.UserData{}, &models.UserRole{})
    if err != nil {
        log.Printf("Error running AutoMigrate: %v", err)
        return err
    }

    log.Println("Database connected and User/Role models migrated successfully.")
    return nil
}