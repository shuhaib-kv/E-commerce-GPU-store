package handler

import (
	"ga/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartHandler struct {
	cartUC *usecase.CartUsecase
}

func NewCartHandler(cuc *usecase.CartUsecase) *CartHandler {
	return &CartHandler{cartUC: cuc}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid user")
		return
	}

	var body struct {
		ProductID string `json:"productid" binding:"required"`
		Quantity  uint   `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	productID, err := primitive.ObjectIDFromHex(body.ProductID)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	result, err := h.cartUC.AddToCart(c.Request.Context(), userID, productID, body.Quantity)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "Added to cart", result)
}

func (h *CartHandler) ViewCart(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid user")
		return
	}

	items, total, err := h.cartUC.ViewCart(c.Request.Context(), userID)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": items, "total_amount": total})
}
