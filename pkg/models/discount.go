package models

import "gorm.io/gorm"

type Discount struct {
	gorm.Model
	DiscountName       string
	DiscountPercentage uint
	ProductId          uint
}
