package controllers

import (
	"errors"
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		Address       uint          `json:"address" binding:"required"`
		Applaywallet  bool          `json:"applaywallet"`
		Coupen        uint          `json:"coupen"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
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

	var cart models.Cart
	if err := database.Db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Cart not found",
			"data":   "null",
		})
		return
	}
	var cartProducts []models.CartProducts
	if err := database.Db.Where("cartid = ?", cart.ID).Find(&cartProducts).Error; err != nil {
		return
	}
	var totalAmount uint
	for _, cp := range cartProducts {
		totalAmount += cp.Quantity * cp.ProductPrice
	}
	if totalAmount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "cart is nill",
			"error":   "add some products to cart",
		})
		return
	}

	if body.Paymentmethod == PaymentMethodCOD {
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
			"data":    fmt.Sprintf("Order placed Expected Delivery Before %s", time.Now().AddDate(0, 0, 12)),
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
	OrderID      string `json:"order_id"`
	Amount       uint   `json:"amount"`
	Deliverydate time.Time
}

func createOrder(cartID uint, userID uint, addressID uint, paymentMethod string) (*CreateOrderResponse, error) {
	if paymentMethod == "cod" {
		var cartProducts []models.CartProducts
		if err := database.Db.Where("cartid = ?", cartID).Find(&cartProducts).Error; err != nil {
			return nil, err
		}

		var totalAmount uint
		for _, cp := range cartProducts {
			totalAmount += cp.ProductPrice
		}

		order := models.Orders{
			UsersID:              userID,
			AddressID:            addressID,
			Orderid:              uuid.New().String(),
			PaymentMethod:        paymentMethod,
			TotalAmount:          totalAmount,
			Status:               true,
			Paymentstatus:        false,
			ExpectedDeliveryDate: time.Now().AddDate(0, 0, 12),
		}
		if err := database.Db.Create(&order).Error; err != nil {
			return nil, err
		}

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

		if err := database.Db.Where("cartid = ?", cartID).Delete(&models.CartProducts{}).Error; err != nil {
			return nil, err
		}

		response := CreateOrderResponse{
			OrderID: order.Orderid,
			Amount:  totalAmount,
		}

		return &response, nil
	} else {
		var cartProducts []models.CartProducts
		if err := database.Db.Where("cartid = ?", cartID).Find(&cartProducts).Error; err != nil {
			return nil, err
		}

		var totalAmount uint
		for _, product := range cartProducts {
			totalAmount += product.ProductPrice
		}

		order := models.Orders{
			UsersID:       userID,
			AddressID:     addressID,
			Orderid:       uuid.New().String(),
			PaymentMethod: paymentMethod,
			TotalAmount:   totalAmount,
			Status:        true,
			Paymentstatus: false,
		}
		if err := database.Db.Create(&order).Error; err != nil {
			return nil, err
		}

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

		if err := database.Db.Where("cartid = ?", cartID).Delete(&models.CartProducts{}).Error; err != nil {
			return nil, err
		}

		response := CreateOrderResponse{
			OrderID:      order.Orderid,
			Amount:       totalAmount,
			Deliverydate: order.ExpectedDeliveryDate,
		}

		return &response, nil
	}
}

type OrderItemResponse struct {
	ProductName string `json:"product_name"`
	Quantity    uint   `json:"quantity"`
	Price       uint   `json:"price"`
}
type OrderResponse struct {
	OrderID       string `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentStatus bool   `json:"payment_status"`
	TotalAmount   uint   `json:"total_amount"`
	Date          time.Time
	Deliverydate  time.Time           `json:"delivery_date"`
	OrderedItems  []OrderItemResponse `json:"ordered_items"`
}

func ListOrders(c *gin.Context) {
	var orders []models.Orders

	// Retrieve all orders
	if err := database.Db.Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving orders"})
		return
	}

	// Retrieve the items associated with each order
	var orderResponses []OrderResponse
	for _, order := range orders {
		var orderedItems []OrderItemResponse
		if err := database.Db.Table("ordereditems").Select("product_name, quantity, price").Where("order_id = ?", order.Orderid).Scan(&orderedItems).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving order items"})
			return
		}

		orderID := order.Orderid
		paymentMethod := order.PaymentMethod
		paymentStatus := order.Paymentstatus
		totalAmount := order.TotalAmount
		deliverydate := order.ExpectedDeliveryDate

		orderResponse := OrderResponse{
			OrderID:       orderID,
			PaymentMethod: paymentMethod,
			PaymentStatus: paymentStatus,
			TotalAmount:   totalAmount,
			Date:          order.CreatedAt,
			OrderedItems:  orderedItems,
			Deliverydate:  deliverydate,
		}
		orderResponses = append(orderResponses, orderResponse)
	}
	c.JSON(http.StatusOK, gin.H{
		"orders":  orderResponses,
		"status":  true,
		"message": "your orders",
	})
}
