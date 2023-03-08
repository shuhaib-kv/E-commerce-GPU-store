package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	User_id uint `json:""`
}
type CartProducts struct {
	gorm.Model
	Cartid       uint   `json:""`
	Productid    uint   `json:""`
	ProductName  string `json:""`
	Quantity     uint
	ProductPrice uint `json:""`
}
