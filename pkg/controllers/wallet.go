package controllers

import (
	"errors"
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletBalance(c *gin.Context) {
	useremail := c.GetString("user")
	fmt.Println(useremail)
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var balance int
	database.Db.Raw("SELECT balance FROM wallets WHERE users_id = ?", UsersID).Scan(&balance)
	c.JSON(200, gin.H{
		"status":  true,
		"Balance": balance,
		"UserID":  UsersID,
	})
}

func WalletInfo(c *gin.Context) {
	useremail := c.GetString("user")
	fmt.Println(useremail)
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var history []models.Wallethistory
	database.Db.Where("wallethistories.users_id = ?", UsersID).Find(&history)

	for _, i := range history {
		c.JSON(200, gin.H{
			"status": true,
			"Date":   i.CreatedAt,
			"Debit":  i.Debit,
			"Credit": i.Credit,
			"UserID": i.UsersID,
		})
	}

}
