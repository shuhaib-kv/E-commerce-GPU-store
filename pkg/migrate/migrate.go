package main

import (
	"ga/pkg/database"
	"ga/pkg/models"
)

func init() {
	database.ConnectDB()
}

func main() {
	database.Db.AutoMigrate(
		&models.Users{},
		&models.Admin{},
		&models.Product{},
		&models.Category{},
		&models.Cart{},
		&models.Orders{},
		&models.Address{},
		&models.Paymentmethod{},
		&models.Ordereditems{},
		&models.CartInfo{},
		&models.RazorPay{},
		&models.Discount{},
		&models.Coupon{},
		&models.Wallet{},
		&models.Wallethistory{},
	)

}
