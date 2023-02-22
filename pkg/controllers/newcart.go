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

func AddToCart(c *gin.Context) {
	useremail := c.GetString("user")
	fmt.Println(useremail)
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}

	var body struct {
		Productid int
		Quantity  int
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
	fmt.Println(product.Stock)

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
	var existingItem models.CartProducts
	if err := database.Db.Where("cartid = ? AND productid = ?", cart.User_id, body.Productid).First(&existingItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"error":  err.Error(),
				"data":   nil,
			})
			return
		}

		if existingItem.ID == uint(body.Productid) {
			newQuantity := existingItem.Quantity + body.Quantity
			if newQuantity > product.Stock {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"error":   "Out of stock",
					"message": "stock will be updated soon",
				})
				return
			}

			if err := database.Db.Model(&existingItem).Update("quantity", newQuantity).Error; err != nil {
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
			}

			c.JSON(http.StatusAccepted, gin.H{
				"status":  true,
				"message": "item present in cart so updated quantity",
				"data":    output,
			})
			return
		}
	}
	// item is in cart, update the quantity

	var cartproduct = models.CartProducts{
		Cartid:       int(cart.ID),
		Productid:    body.Productid,
		ProductName:  product.Name,
		Quantity:     body.Quantity,
		ProductPrice: product.Price,
	}

	if err := database.Db.Create(&cartproduct).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  err,
			"data":   "null",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  true,
		"data":    cartproduct,
		"message": "Added to cart",
	})
}
func ViewCart(c gin.Context) {

}
