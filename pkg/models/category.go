package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	ID   int    `json:"id" gorm:"primaryKey;unique"  `
	Name string `json:"name"`
}
