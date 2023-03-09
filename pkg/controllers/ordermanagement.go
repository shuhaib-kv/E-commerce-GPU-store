package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ViewOrders(c *gin.Context) {
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		pageNum = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = 10
	}
	var ps bool
	searchOrderID := c.Query("orderID")
	searchStatus := c.Query("status")
	searchPaymentMethod := c.Query("paymentMethod")
	searchPaymentStatus := c.Query("paymentStatus")
	if searchPaymentStatus == "true" {
		ps = true
	} else if searchPaymentStatus == "false" {
		ps = false
	}
	offset := (pageNum - 1) * pageSize
	limit := pageSize

	var orders []models.Orders
	dbQuery := database.Db.Offset(offset).Limit(limit)

	if searchOrderID != "" {
		dbQuery = dbQuery.Where("orderid = ?", searchOrderID)
	}
	if searchStatus != "" {
		dbQuery = dbQuery.Where("status = ?", searchStatus)
	}
	if searchPaymentMethod != "" {
		dbQuery = dbQuery.Where("payment_method = ?", ps)
	}
	if searchPaymentStatus != "" {
		dbQuery = dbQuery.Where("paymentstatus = ?", searchPaymentStatus)
	}

	dbQuery.Find(&orders)

	var count int64
	dbQuery.Model(&models.Orders{}).Count(&count)

	if len(orders) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No orders found"})
		return
	}

	var ordersList []interface{}

	for _, order := range orders {
		orderDetails := map[string]interface{}{
			"userid":        order.UsersID,
			"addressid":     order.AddressID,
			"orderid":       order.Orderid,
			"paymentmethod": order.PaymentMethod,
			"totalamount":   order.TotalAmount,
			"paymentstatus": order.Paymentstatus,
			"status":        order.Status,
		}
		ordersList = append(ordersList, orderDetails)
	}

	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{"status": true, "data": ordersList, "count": count, "totalPages": totalPages})
}

func EditOrder(c *gin.Context) {
	var body struct {
		Oderid        string `json:"orderid"`
		Paymentstatus bool   `json:"paymentstatus"`
		Orderstatus   bool   `json:"orderstatus"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "binding json faild",
			"data":    "error ",
		})
		return
	}
	order := models.Orders{}
	result := database.Db.Where("orderid = ?", body.Oderid).First(&order)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "order not found",
			"data":    "error",
		})
		return
	}
	order.Paymentstatus = body.Paymentstatus
	order.Status = body.Orderstatus
	result = database.Db.Save(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to update order",
			"data":    "error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "order updated successfully",
		"data":    order,
	})
}
