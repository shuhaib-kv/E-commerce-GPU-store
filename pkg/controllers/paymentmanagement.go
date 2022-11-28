package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"

	"github.com/gin-gonic/gin"
)

func AddPaymentMethod(c *gin.Context) {
	var body struct {
		Payment_Method string `json:"paymentmethod"`
	}
	c.Bind(&body)
	pay := models.Paymentmethod{
		Payment_Method: body.Payment_Method,
	}
	createpaymentmethod := database.Db.Create(&pay)
	if createpaymentmethod.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "method added",
	})
}
