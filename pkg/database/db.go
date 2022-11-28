package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDB() {
	var err error
	// dsn := os.Getenv("DB")
	Db, err = gorm.Open(postgres.Open("host=localhost user=postgres  dbname=gpu_ecom port=5432 "), &gorm.Config{})
	if err != nil {
		fmt.Println("error", err)
	}

}
