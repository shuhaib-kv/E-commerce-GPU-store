package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//		func ViewOrders(c *gin.Context) {
//			var order []models.Orders
//			database.Db.Find(&order)
//			for _, i := range order {
//				c.JSON(200, gin.H{
//					"id":             i.OrderID,
//					"user id":        i.UsersID,
//					"price":          i.Total_Amount,
//					"Adressid":       i.AddressID,
//					"order status":   i.OrderStatus,
//					"payment status": i.PaymentStatus,
//				})
//	}
func ViewOrders(c *gin.Context) {
	// Get the page number and page size from the query parameters
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		pageNum = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = 10
	}
	var ps bool

	// Get the search terms from the query parameters
	searchOrderID := c.Query("orderID")
	searchStatus := c.Query("status")
	searchPaymentMethod := c.Query("paymentMethod")
	searchPaymentStatus := c.Query("paymentStatus")
	if searchPaymentStatus == "true" {
		ps = true
	} else if searchPaymentStatus == "false" {
		ps = false
	}
	// Calculate the offset and limit for the query
	offset := (pageNum - 1) * pageSize
	limit := pageSize

	// Query the database to get the orders
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

	// Get the total number of orders for pagination
	var count int64
	dbQuery.Model(&models.Orders{}).Count(&count)

	if len(orders) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No orders found"})
		return
	}

	// Create a slice to hold the response
	var ordersList []interface{}

	// Loop through the orders and add the necessary fields to the response slice
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

	// Calculate the total number of pages for pagination
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))

	// Return the response
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": ordersList, "count": count, "totalPages": totalPages})
}

// func EditOrder(c *gin.Context) {
// 	id := c.Param("id")
// 	var body struct {
// 		PaymentStatus string
// 		OrderStatus   string
// 	}
// 	c.Bind(&body)
// 	var order []models.Orders
// 	result := database.Db.First(&order, id)
// 	database.Db.Model(&order).Updates(models.Orders{
// 		PaymentStatus: body.PaymentStatus,
// 		OrderStatus:   body.OrderStatus,
// 	})
// 	if result != nil {
// 		c.JSON(200, gin.H{
// 			"message": " can't find product",
// 		})
// 	} else {
// 		c.JSON(200, gin.H{
// 			"message": "order updated",
// 		})
// 	}
// }
// fp
