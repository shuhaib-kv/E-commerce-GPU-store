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

type PaymentMethod string

const (
	PaymentMethodCOD      PaymentMethod = "cod"
	PaymentMethodRazorpay PaymentMethod = "razorpay"
)

func OrderCart(c *gin.Context) {
	var body struct {
		Paymentmethod PaymentMethod `json:"payment_method" binding:"required,oneof=cod razorpay"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID
	useremail := c.GetString("user")
	var userID uint
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&userID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "User not found",
		})
		return
	}

	// Find cart for the user
	var cart models.Cart
	if err := database.Db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Cart not found",
			"data":   "null",
		})
		return
	}

	// Find products in the cart
	var cartProducts []models.CartProducts
	database.Db.Where("cartid = ?", cart.ID).Find(&cartProducts)

	// Calculate total amount
	var totalAmount uint
	for _, product := range cartProducts {
		totalAmount += product.ProductPrice
	}

	// Create order in database
	order := models.Orders{
		UsersID:       userID,
		PaymentMethod: string(body.Paymentmethod),
		TotalAmount:   totalAmount,
		Status:        123,
	}
	if err := database.Db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create order",
			"error":   err.Error(),
		})
		return
	}

	// Clear cart
	database.Db.Where("cartid = ?", cart.ID).Delete(&models.CartProducts{})
	database.Db.Delete(&cart)

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"message":  fmt.Sprintf("Order placed successfully with payment method %s", body.Paymentmethod),
		"order_id": order,
	})

}
