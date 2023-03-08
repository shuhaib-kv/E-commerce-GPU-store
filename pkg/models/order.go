package models

import (
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	UsersID              uint
	AddressID            uint
	Orderid              string
	PaymentMethod        string
	TotalAmount          uint
	Status               bool
	Paymentstatus        bool
	ExpectedDeliveryDate time.Time
}

type Ordereditems struct {
	gorm.Model
	OrderID     string
	ProductID   uint
	ProductName string
	Quantity    uint
	Price       uint
}
