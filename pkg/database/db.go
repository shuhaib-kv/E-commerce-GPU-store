package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDB() {
	var err error

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error connecting to database")

	}
	// Db, err = gorm.Open(postgres.Open("host=localhost user=soib password=soib  dbname=gpu_ecom port=5432 "), &gorm.Config{})
	// if err != nil {
	// 	fmt.Println("error", err)
	// }

}
