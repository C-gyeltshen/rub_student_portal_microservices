package models

import (
    "time"
    "gorm.io/gorm"
)

type UserData struct {
    gorm.Model
    // Removed duplicate ID since gorm.Model includes it

    First_name  string `json:"first name"`
    Second_name string `json:"second name"`
    Email       string `json:"email"`

    // 1. Foreign Key: This column will store the ID of the related role.
    UserRoleID uint

    // 2. Struct Relationship: This is the actual role object GORM loads.
    // The tag `gorm:"foreignKey:UserRoleID"` explicitly links it to the FK column.
    Role UserRole `gorm:"foreignKey:UserRoleID"`

    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}