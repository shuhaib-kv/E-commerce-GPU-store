package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	Product_ID    int`json:""`
	User_id       int`json:""`
	Product_Name  string`json:""`
	Brand_Name    string`json:""`
	Description   string`json:""`
	Quantity      int`json:""`
	Price         int`json:""`
	DiscountPrice int`json:""`
	Total         int`json:""`
}
type CartInfo struct {
	gorm.Model
	UsersId      int`json:""`
	ProductsID   int`json:""`
	Discount     int`json:""`
	ProductName  string`json:""`
	BrandName    string`json:""`
	ProductPrice int`json:""`
}
type Checkoutinfo struct {
	gorm.Model
	UsersID        int`json:""`
	OrderID        string`json:""`
	Discount       int`json:""`
	CouponDiscount int`json:""`
	CouponCode     string`json:""`
	TotalMrp       int`json:""`
	Total          int`json:""`
}
