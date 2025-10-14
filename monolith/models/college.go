package models

import (

	"gorm.io/gorm"
)

type College struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}