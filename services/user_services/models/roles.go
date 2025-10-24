package models

import (
    "time"

    "gorm.io/gorm"
)

type UserRole struct {
    gorm.Model
    // Removed duplicate ID since gorm.Model includes it

    Name        string         `json:"name"`
    Description string         `json:"description"`
    
    // Optional: You can add a `User` field here to traverse the relationship
    // from Role back to User, though it's not strictly necessary for the link.
    // User UserData 
    
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}