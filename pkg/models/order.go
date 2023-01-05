package models

import "gorm.io/gorm"

type Orders struct {
	UsersID        int
	AddressID      int
	OrderID        string
	Ordertype      string
	Discount       int
	CouponDiscount int
	CouponCode     string
	Payment_Method string
	Total_Amount   int
	PaymentStatus  string
	OrderStatus    string
}

type Ordereditems struct {
	gorm.Model
	UsersID        int
	ProductsID     int
	Order_ID       string
	Product_Name   string
	Price          int
	CouponDiscount int
	Discount       int
	PaymentStatus  string
	OrderStatus    string
	Payment_Method string
	Totalamount    int
	//AmountPaid     int
}
