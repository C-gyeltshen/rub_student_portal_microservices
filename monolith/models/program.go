package models

import "gorm.io/gorm"

type Program struct {
    gorm.Model
    Name  string `json:"name"`
    Description string `json:"description"`
}

