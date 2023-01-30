package models

import "gorm.io/gorm"

type Orders struct {
	UsersID        int`json:""`
	AddressID      int`json:""`
	OrderID        string`json:""`
	Ordertype      string`json:""`
	Discount       int`json:""`
	CouponDiscount int`json:""`
	CouponCode     string`json:""`
	Payment_Method string`json:""`
	Total_Amount   int`json:""`
	PaymentStatus  string`json:""`
	OrderStatus    string`json:""`
}

type Ordereditems struct {
	gorm.Model
	UsersID        int`json:""`
	ProductsID     int`json:""`
	Order_ID       string`json:""`
	Product_Name   string`json:""`
	Price          int`json:""`
	CouponDiscount int`json:""`
	Discount       int`json:""`
	PaymentStatus  string`json:""`
	OrderStatus    string`json:""`
	Payment_Method string`json:""`
	Totalamount    int`json:""`
	//AmountPaid     int
}
