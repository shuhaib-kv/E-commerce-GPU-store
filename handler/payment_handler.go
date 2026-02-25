package handler

import (
	"ga/config"
	"ga/domain"
	"ga/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentHandler struct {
	paymentUC *usecase.PaymentUsecase
	cfg       *config.Config
}

func NewPaymentHandler(puc *usecase.PaymentUsecase, cfg *config.Config) *PaymentHandler {
	return &PaymentHandler{paymentUC: puc, cfg: cfg}
}

func (h *PaymentHandler) RazorPay(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user"})
		return
	}

	order, err := h.paymentUC.GetPendingOrder(c.Request.Context(), userID)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "No pending order"})
		return
	}

	client := razorpay.NewClient(h.cfg.RazorpayKey, h.cfg.RazorpaySecret)

	data := map[string]interface{}{
		"amount":   order.TotalAmount * 100,
		"currency": "INR",
		"receipt":  order.OrderID,
	}
	rzpOrder, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to create payment"})
		return
	}

	c.HTML(http.StatusOK, "app.html", gin.H{
		"userid":      userID.Hex(),
		"totalamount": order.TotalAmount,
		"orderid":     rzpOrder["id"],
		"email":       c.GetString("user_email"),
		"phone":       "",
		"razorpaykey": h.cfg.RazorpayKey,
	})
}

func (h *PaymentHandler) RazorpaySuccess(c *gin.Context) {
	userIDStr := c.Query("user_id")
	paymentID := c.Query("payment_id")
	orderID := c.Query("order_id")
	signature := c.Query("signature")
	amountStr := c.Query("id")

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user"})
		return
	}

	amount, _ := strconv.Atoi(amountStr)

	payment := &domain.RazorPayment{
		UserID:          userID,
		RazorPaymentID:  paymentID,
		RazorPayOrderID: orderID,
		Signature:       signature,
		AmountPaid:      uint(amount),
	}

	if err := h.paymentUC.SavePayment(c.Request.Context(), payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	// Mark order as paid - find the order by looking up the pending order
	h.paymentUC.MarkOrderPaid(c.Request.Context(), orderID)

	c.Redirect(http.StatusFound, "/success")
}

func (h *PaymentHandler) Success(c *gin.Context) {
	c.HTML(http.StatusOK, "success.html", nil)
}
