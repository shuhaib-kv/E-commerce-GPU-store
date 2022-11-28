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
	Db, err = gorm.Open(postgres.Open("host=localhost user=soib password=soib dbname=gpu_ecom port=5432 sslmode=prefer"), &gorm.Config{})
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Printf("Database Type ='%s'\n Connected to database sussesfully!", Db.Name())

}
