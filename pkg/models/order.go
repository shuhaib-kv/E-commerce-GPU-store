package models

import "gorm.io/gorm"

type Orders struct {
	UsersID       uint
	AddressID     uint
	Orderid       string
	PaymentMethod string
	TotalAmount   uint
	Status        bool
	Paymentstatus bool
}

type Ordereditems struct {
	gorm.Model
	OrderID     string `json:"order_id"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    uint   `json:"quantity"`
	Price       uint   `json:"price"`
}
