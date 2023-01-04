package models

import "gorm.io/gorm"

type Wallet struct {
	UsersID int
	Balance int
}
type Wallethistory struct {
	gorm.Model
	UsersID uint
	Credit  int
	Debit   int
}
