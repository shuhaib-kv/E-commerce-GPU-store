package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPaymentMethod(c *gin.Context) {
	var request struct {
		PaymentMethod string `json:"payment_method"`
	}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"error":   err,
			"message": "Invalid request body",
		})
		return
	}
	paymentMethod := models.Paymentmethod{
		Payment_Method: request.PaymentMethod,
	}
	result := database.Db.Create(&paymentMethod)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"error":   result.Error,
			"message": "Failed to add payment method",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Payment method added successfully",
		"data":    paymentMethod,
	})
}
