package database

import (
	"fmt"
	"ga/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDB() {
	//initialized database
	var err error
	Db, err = gorm.Open(postgres.Open("host=localhost user=postgres password=soib  dbname=gpu port=5432 "), &gorm.Config{})
	if err != nil {
		fmt.Println("error", err)
		return
	}
	Db.AutoMigrate(
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
