package controllers

import (
	"errors"
	"ga/pkg/database"
	"ga/pkg/models"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var letters = []rune("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func OrderIdGeneration(value int) string {
	b := make([]rune, value)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func OrderInfo(c *gin.Context) {

	var Orders []models.Orders
	database.Db.Find(&Orders)
	for _, i := range Orders {
		c.JSON(200, gin.H{
			"status":         true,
			"UserId":         i.UsersID,
			"OrderID":        i.OrderID,
			"Discount":       i.Discount,
			"CouponDiscount": i.CouponDiscount,
			"CouponCode":     i.CouponCode,
			"PaymentMethod":  i.Payment_Method,
			"TotalAmount":    i.Total_Amount,
		})
	}

}

func OrderedItems(c *gin.Context) {
	// Fetching user id from jwt
	useremail := c.GetString("user")
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var items []models.Orders
	database.Db.Where("users_id = ?", UsersID).Find(&items)

	for _, i := range items {
		c.JSON(200, gin.H{
			"status":      true,
			"id":          i.OrderID,
			"Amount_Paid": i.Total_Amount,

			"Discount":        i.Discount,
			"Coupon_Discount": i.CouponDiscount,
			"Order Status":    i.OrderStatus,
		})
	}
}

func CancelOrder(c *gin.Context) {
	useremail := c.GetString("user")
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})
		var items models.Orders
		var updateStatus string = "CANCELLED"
		id := c.Query("orderid")

		database.Db.First(&items, id)
		if items.OrderStatus == updateStatus {
			c.JSON(400, gin.H{
				"status":  false,
				"message": "Order already Cancelled",
			})
			return
		}
		database.Db.Model(&items).Where("id=?", id).Update("order_status", updateStatus)

		var price int
		database.Db.Raw("SELECT price FROM ordereditems WHERE id = ?", id).Scan(&price)

		var balance int
		database.Db.Raw("SELECT balance FROM wallets WHERE users_id = ?", UsersID).Scan(&balance)
		newBalance := balance + price

		if items.Payment_Method == "COD" {
			c.JSON(200, gin.H{
				"status":  true,
				"message": "Order Cancelled",
			})
			return
		}

		WalletHistory := models.Wallethistory{UsersID: uint(UsersID), Debit: 0, Credit: price}
		database.Db.Create(&WalletHistory)

		var totalAmount int
		database.Db.Raw("SELECT total_amaount FROM orders WHERE users_id = ?", UsersID).Scan(&totalAmount)
		Ntotal := totalAmount - balance
		// Updating wallet on order cancellation
		database.Db.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
		database.Db.Model(&models.Orders{}).Where("users_id = ?", UsersID).Update("total_amount", Ntotal)
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Order Cancelled",
		})
	}
}
func CreatesOrderId() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "OID" + id
	return orderID
}
