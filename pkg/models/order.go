package models

import "gorm.io/gorm"

type Orders struct {
	UsersID        uint
	AddressID      uint
	OrderID        string
	Ordertype      string
	Discount       uint
	CouponDiscount uint
	CouponCode     string
	Payment_Method string
	Total_Amount   uint
	PaymentStatus  string
	OrderStatus    string
}

type Ordereditems struct {
	gorm.Model
	UsersID        uint
	ProductsID     uint
	Order_ID       string
	Product_Name   string
	Price          uint
	CouponDiscount uint
	Discount       uint
	PaymentStatus  string
	OrderStatus    string
	Payment_Method string
	Totalamount    uint
	//AmountPaid     int
}
