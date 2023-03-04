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
	Db, err = gorm.Open(postgres.Open("host=localhost user=postgres password=soib  dbname=shuhaib port=5432 "), &gorm.Config{})
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
		&models.CartProducts{},
		&models.Orders{},
		&models.Address{},
		&models.Paymentmethod{},
		&models.Ordereditems{},
		&models.RazorPay{},
		&models.Discount{},
		&models.Coupon{},
		&models.Wallet{},
		&models.Wallethistory{},
	)
}
