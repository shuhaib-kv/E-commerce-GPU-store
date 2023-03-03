package controllers

import (
	"errors"
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentMethodCOD      PaymentMethod = "cod"
	PaymentMethodRazorpay PaymentMethod = "razorpay"
)

// func OrderCart(c *gin.Context) {
// var body struct {
// 	Paymentmethod PaymentMethod `json:"payment_method" binding:"required,oneof=cod razorpay"`
// }

// if err := c.ShouldBindJSON(&body); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{
// 		"status":  false,
// 		"message": "Invalid request body",
// 		"error":   err.Error(),
// 	})
// 	return
// }

// // Get user ID
// useremail := c.GetString("user")
// var userID uint
// err := database.Db.Raw("select id from users where email=?", useremail).Scan(&userID)
// if errors.Is(err.Error, gorm.ErrRecordNotFound) {
// 	c.JSON(http.StatusNotFound, gin.H{
// 		"status":  false,
// 		"message": "User not found",
// 	})
// 	return
// }

// // Find cart for the user
// var cart models.Cart
// if err := database.Db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
// 	c.JSON(http.StatusNotFound, gin.H{
// 		"status": false,
// 		"error":  "Cart not found",
// 		"data":   "null",
// 	})
// 	return
// }

// // Find products in the cart
// var cartProducts []models.CartProducts
// database.Db.Where("cartid = ?", cart.ID).Find(&cartProducts)

// // Calculate total amount
// var totalAmount uint
// for _, product := range cartProducts {
// 	totalAmount += product.ProductPrice
// }

// // Create order in database
// order := models.Orders{
// 	UsersID:       userID,
// 	PaymentMethod: string(body.Paymentmethod),
// 	TotalAmount:   totalAmount,
// 	Status:        123,
// }
// if err := database.Db.Create(&order).Error; err != nil {
// 	c.JSON(http.StatusInternalServerError, gin.H{
// 		"status":  false,
// 		"message": "Failed to create order",
// 		"error":   err.Error(),
// 	})
// 	return
// }

// // Clear cart
// database.Db.Where("cartid = ?", cart.ID).Delete(&models.CartProducts{})
// database.Db.Delete(&cart)

// c.JSON(http.StatusOK, gin.H{
// 	"status":   true,
// 	"message":  fmt.Sprintf("Order placed successfully with payment method %s", body.Paymentmethod),
// 	"order_id": order,
// 	})

// }
func OrderCart(c *gin.Context) {
	var body struct {
		Paymentmethod PaymentMethod `json:"payment_method" binding:"required,oneof=cod razorpay"`
		Address       uint          `json:"address"`
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
	if body.Paymentmethod == PaymentMethodCOD {
		// Create the order from the cart
		if _, err := createOrder(cart.ID, userID, body.Address, string(body.Paymentmethod)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to create order",
				"error":   err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": fmt.Sprintf("Order placed successfully with payment method %s", body.Paymentmethod),
			"data":    "",
		})
		return
	}
	if body.Paymentmethod == PaymentMethodRazorpay {
		if _, err := createOrder(cart.ID, userID, body.Address, string(body.Paymentmethod)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to create order",
				"error":   err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": fmt.Sprintf("Order payment method %s", body.Paymentmethod),
			"data":    "Go to Rzorpay to complete payment",
		})
		return
	}

}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	Amount  uint   `json:"amount"`
}

func createOrder(cartID uint, userID uint, addressID uint, paymentMethod string) (*CreateOrderResponse, error) {
	// Get the products in the cart
	var cartProducts []models.CartProducts
	if err := database.Db.Where("cartid = ?", cartID).Find(&cartProducts).Error; err != nil {
		return nil, err
	}

	// Calculate the total amount of the order
	var totalAmount uint
	for _, cp := range cartProducts {
		totalAmount += cp.Quantity * cp.ProductPrice
	}
	// Create the order
	order := models.Orders{
		UsersID:       userID,
		AddressID:     addressID,
		Orderid:       uuid.New().String(), // generate a unique order ID using UUID
		PaymentMethod: paymentMethod,
		TotalAmount:   totalAmount,
		Status:        true,
		Paymentstatus: false,
	}
	if err := database.Db.Create(&order).Error; err != nil {
		return nil, err
	}

	// Add the products in the cart to the order
	for _, cp := range cartProducts {
		orderedItem := models.Ordereditems{
			OrderID:     order.Orderid,
			ProductID:   cp.Productid,
			ProductName: cp.ProductName,
			Quantity:    cp.Quantity,
			Price:       cp.ProductPrice,
		}
		if err := database.Db.Create(&orderedItem).Error; err != nil {
			return nil, err
		}
	}

	// Delete the products in the cart
	if err := database.Db.Where("cartid = ?", cartID).Delete(&models.CartProducts{}).Error; err != nil {
		return nil, err
	}

	// Create the response
	response := CreateOrderResponse{
		OrderID: order.Orderid,
		Amount:  totalAmount,
	}

	return &response, nil
}
