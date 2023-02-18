package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	Product_ID int
	User_id    int
	Quantity   int
	Price      int
	Total      int
}

type Cartinfo struct {
	gorm.Model
	UsersID      int
	OrderID      string
	ProductPrice int
	Total        int
}
