package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	UsersID uint
	Balance uint
}
type Wallethistory struct {
	gorm.Model
	UsersID uint
	Credit  uint
	Debit   uint
}
