package models

import "gorm.io/gorm"

type Wallet struct {
	UsersID int`json:""`
	Balance int`json:""`
}
type Wallethistory struct {
	gorm.Model
	UsersID uint`json:""`
	Credit  int`json:""`
	Debit   int`json:""`
}
