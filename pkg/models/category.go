package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	ID   uint   `json:"id" gorm:"primaryKey;unique"  `
	Name string `json:"name"`
}
