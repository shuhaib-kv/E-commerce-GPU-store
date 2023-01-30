package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey;unique"  `
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
