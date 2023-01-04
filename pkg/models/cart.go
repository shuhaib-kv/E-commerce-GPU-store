package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	Product_ID    int
	User_id       int
	Product_Name  string
	Brand_Name    string
	Description   string
	Quantity      int
	Price         int
	DiscountPrice int
	Total         int
}
type CartInfo struct {
	gorm.Model
	UsersId      int
	ProductsID   int
	Discount     int
	ProductName  string
	BrandName    string
	ProductPrice int
}
type Checkoutinfo struct {
	gorm.Model
	UsersID        int
	OrderID        string
	Discount       int
	CouponDiscount int
	CouponCode     string
	TotalMrp       int
	Total          int
}
