package controllers

import (
	"errors"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddToCart(c *gin.Context) {
	useremail := c.GetString("user")
	var UsersID uint
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var body struct {
		Productid uint
		Quantity  uint
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err,
		})
		return
	}
	var product models.Product
	if err := database.Db.Where("id = ?", body.Productid).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "can't find product",
			"data":   "null",
		})
		return
	}

	var cart models.Cart
	if err := database.Db.Find(&cart, UsersID).Scan(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Cart not found",
			"data":   "null",
		})
		return
	}
	if body.Quantity > product.Stock {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"error":   "Out of stock",
			"message": "stock will be updated soon",
		})
		return
	}
	//
	var existingItems []models.CartProducts
	database.Db.Where("cartid = ?", cart.User_id).Find(&existingItems)
	for _, existingItem := range existingItems {
		if existingItem.Productid == body.Productid {
			newQuantity := existingItem.Quantity + body.Quantity
			newPrice := product.Price * newQuantity
			if newQuantity > product.Stock {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"error":   "Out of stock",
					"message": "stock will be updated soon",
				})
				return
			}
			if err := database.Db.Model(&existingItem).Updates(models.CartProducts{
				Quantity:     newQuantity,
				ProductPrice: newPrice,
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": false,
					"error":  err.Error(),
					"data":   nil,
				})
				return
			}

			output := make([]map[string]interface{}, 1)
			output[0] = map[string]interface{}{
				"cartid":      existingItem.Cartid,
				"productid":   existingItem.Productid,
				"productname": existingItem.ProductName,
				"quantity":    newQuantity,
				"totalprice":  existingItem.ProductPrice,
			}

			c.JSON(http.StatusAccepted, gin.H{
				"status":  true,
				"message": "item present in cart so updated quantity",
				"data":    output,
			})
			return
		}
	}
	var cartproduct = models.CartProducts{
		Cartid:       uint(cart.ID),
		Productid:    body.Productid,
		ProductName:  product.Name,
		Quantity:     body.Quantity,
		ProductPrice: product.Price * body.Quantity,
	}

	if err := database.Db.Create(&cartproduct).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  err,
			"data":   "null",
		})
		return
	}
	output := make([]map[string]interface{}, 1)
	output[0] = map[string]interface{}{
		"cartid":      cartproduct.Cartid,
		"productid":   cartproduct.Productid,
		"productname": cartproduct.ProductName,
		"quantity":    body.Quantity,
		"totalprice":  cartproduct.ProductPrice,
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  true,
		"data":    output,
		"message": "Added to cart",
	})

}

func ViewCart(c *gin.Context) {
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
	database.Db.Where("cartid = ?", cart.ID).Find(&cartProducts)

	var totalAmount uint
	for _, product := range cartProducts {
		totalAmount += product.ProductPrice
	}

	var output struct {
		CartID       uint                     `json:"cart_id"`
		TotalAmount  uint                     `json:"total_amount"`
		CartProducts []map[string]interface{} `json:"cart_products"`
	}
	output.CartID = cart.ID
	output.TotalAmount = totalAmount
	for _, product := range cartProducts {
		output.CartProducts = append(output.CartProducts, map[string]interface{}{
			"name":  product.ProductName,
			"price": product.ProductPrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "items in your cart",
		"data":    output,
	})
}
