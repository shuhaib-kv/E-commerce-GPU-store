package handler

import (
	"ga/usecase"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	orderUC *usecase.OrderUsecase
}

func NewOrderHandler(ouc *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUC: ouc}
}

func (h *OrderHandler) OrderCart(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user"})
		return
	}

	var body struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
		AddressID     string `json:"address_id" binding:"required"`
		ApplyWallet   bool   `json:"apply_wallet"`
		CouponCode    string `json:"coupon_code"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	addressID, err := primitive.ObjectIDFromHex(body.AddressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid address ID"})
		return
	}

	orderID, amount, err := h.orderUC.CreateOrder(c.Request.Context(), userID, addressID, body.PaymentMethod, body.CouponCode, body.ApplyWallet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true, "message": "Order placed",
		"order_id": orderID, "total_amount": amount,
	})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user"})
		return
	}

	orders, err := h.orderUC.ListOrders(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": orders})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user"})
		return
	}

	var body struct {
		OrderID string `json:"orderid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	if err := h.orderUC.CancelOrder(c.Request.Context(), userID, body.OrderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Order cancelled"})
}

// Admin
func (h *OrderHandler) ViewOrders(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)

	filter := map[string]interface{}{}
	if orderID := c.Query("orderID"); orderID != "" {
		filter["order_id"] = orderID
	}
	if status := c.Query("status"); status != "" {
		filter["status"] = status == "true"
	}
	if pm := c.Query("paymentMethod"); pm != "" {
		filter["payment_method"] = pm
	}
	if ps := c.Query("paymentStatus"); ps != "" {
		filter["payment_status"] = ps == "true"
	}

	orders, total, err := h.orderUC.ListAllOrders(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, gin.H{
		"status": true, "data": orders, "currentPage": page,
		"pageSize": pageSize, "totalPages": totalPages, "totalItems": total,
	})
}

func (h *OrderHandler) EditOrder(c *gin.Context) {
	var body struct {
		OrderID       string `json:"orderid" binding:"required"`
		Status        string `json:"status"`
		PaymentStatus string `json:"paymentstatus"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	status, _ := strconv.ParseBool(body.Status)
	paymentStatus, _ := strconv.ParseBool(body.PaymentStatus)

	if err := h.orderUC.EditOrder(c.Request.Context(), body.OrderID, status, paymentStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Order updated"})
}
