package controllers

import (
	"errors"
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"strconv"
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
		Coupen        string        `json:"coupen"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
	userID, _ := strconv.ParseUint(c.GetString("id"), 10, 32)

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
		if _, err := createOrder(cart.ID, uint(userID), body.Address, string(body.Paymentmethod), body.Coupen, body.Applaywallet); err != nil {
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
		if _, err := createOrder(cart.ID, uint(userID), body.Address, string(body.Paymentmethod), body.Coupen, body.Applaywallet); err != nil {
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

func createOrder(cartID uint, userID uint, addressID uint, paymentMethod string, coupen string, applaywallet bool) (*CreateOrderResponse, error) {
	if paymentMethod == "cod" {
		var cartProducts []models.CartProducts
		if err := database.Db.Where("cartid = ?", cartID).Find(&cartProducts).Error; err != nil {
			return nil, err
		}

		var totalAmount uint
		for _, cp := range cartProducts {
			product := models.Product{}
			if err := database.Db.First(&product, cp.Productid).Error; err != nil {
				return nil, err
			}
			if product.Discount != 0 {
				discount := models.Discount{}
				if err := database.Db.First(&discount, product.Discount).Error; err != nil {
					return nil, err
				}
				productPriceWithDiscount := (product.Price * (100 - discount.DiscountPercentage)) / 100
				totalAmount += (productPriceWithDiscount * cp.Quantity)
			} else {
				totalAmount += (product.Price * cp.Quantity)
			}
		}

		var coupon models.Coupon
		if coupen != "" {
			if err := database.Db.Where("coupon_code = ?", coupen).First(&coupon).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("invalid coupon code")
				}
				return nil, err
			}
			if time.Now().After(coupon.ExpiryDate) {
				return nil, errors.New("coupon has expired")
			}
			discount := (coupon.CouponPercentage * totalAmount) / 100
			totalAmount -= discount
		}

		var balance uint
		if applaywallet {
			var wallet models.Wallet
			if err := database.Db.Where("users_id = ?", userID).First(&wallet).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("user wallet not found")
				}
				return nil, err
			}
			if wallet.Balance < totalAmount {
				return nil, errors.New("insufficient wallet balance")
			}
			balance = wallet.Balance - totalAmount
			if err := database.Db.Model(&wallet).Update("balance", balance).Error; err != nil {
				return nil, err
			}
			var wallethistory models.Wallethistory
			wallethistory.UsersID = userID
			wallethistory.Credit = 0
			wallethistory.Debit = totalAmount
			if err := database.Db.Create(&wallethistory).Error; err != nil {
				return nil, err
			}
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

		var coupon models.Coupon
		if err := database.Db.Where("coupon_code = ?", coupen).First(&coupon).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("invalid coupon code")
			}
			return nil, err
		}
		if time.Now().After(coupon.ExpiryDate) {
			return nil, errors.New("coupon has expired")
		}

		discount := (coupon.CouponPercentage * totalAmount) / 100
		totalAmount -= discount

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

	if err := database.Db.Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving orders"})
		return
	}

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
