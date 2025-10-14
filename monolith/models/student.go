package models

import "gorm.io/gorm"

type Student struct {
    gorm.Model
	Studnet_id int `json:"student_id" gorm:"unique"`
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
}