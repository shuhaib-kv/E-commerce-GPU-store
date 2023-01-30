package models

import "gorm.io/gorm"

type Paymentmethod struct {
	gorm.Model
	Payment_Method string`json:""`
}
