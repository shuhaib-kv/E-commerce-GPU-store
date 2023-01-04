package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"

	"github.com/gin-gonic/gin"
)

func ViewOrders(c *gin.Context) {
	var order []models.Orders
	database.Db.Find(&order)
	for _, i := range order {
		c.JSON(200, gin.H{
			"id":             i.OrderID,
			"user id":        i.UsersID,
			"price":          i.Total_Amount,
			"Adressid":       i.AddressID,
			"order status":   i.OrderStatus,
			"payment status": i.PaymentStatus,
		})
	}

}
func EditOrder(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		PaymentStatus string
		OrderStatus   string
	}
	c.Bind(&body)
	var order []models.Orders
	result := database.Db.First(&order, id)
	database.Db.Model(&order).Updates(models.Orders{
		PaymentStatus: body.PaymentStatus,
		OrderStatus:   body.OrderStatus,
	})
	if result != nil {
		c.JSON(200, gin.H{
			"message": " can't find product",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "order updated",
		})
	}
}
