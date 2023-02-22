package models

import "gorm.io/gorm"

type Wallet struct {
	UsersID uint
	Balance uint
}
type Wallethistory struct {
	gorm.Model
	UsersID uint
	Credit  uint
	Debit   uint
}
