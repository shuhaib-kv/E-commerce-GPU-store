package models

import "gorm.io/gorm"

type Discount struct {
	gorm.Model
	DiscountName       string`json:""`
	DiscountPercentage int`json:""`
	ProductId          int`json:""`
}
