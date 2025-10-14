package models

import "gorm.io/gorm"

type Bank struct {
    gorm.Model
	Account_number int `json:"Account number" gorm:"unique"`
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
}