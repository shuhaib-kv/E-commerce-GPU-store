package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model

	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type In struct {
	gorm.Model
	Name string
}
