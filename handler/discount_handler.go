package handler

import (
	"ga/usecase"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiscountHandler struct {
	discountUC *usecase.DiscountUsecase
}

func NewDiscountHandler(duc *usecase.DiscountUsecase) *DiscountHandler {
	return &DiscountHandler{discountUC: duc}
}

func (h *DiscountHandler) AddDiscount(c *gin.Context) {
	var body struct {
		DiscountName       string `json:"discountname" binding:"required"`
		DiscountPercentage uint   `json:"discountpercentage" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	discount, err := h.discountUC.AddDiscount(c.Request.Context(), body.DiscountName, body.DiscountPercentage)
	if err != nil {
		respondInternalError(c, err, "AddDiscount")
		return
	}
	respondSuccess(c, http.StatusCreated, "Discount created", discount)
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	var body struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := primitive.ObjectIDFromHex(body.ID)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.discountUC.DeleteDiscount(c.Request.Context(), id); err != nil {
		respondInternalError(c, err, "DeleteDiscount")
		return
	}
	respondSuccess(c, http.StatusOK, "Discount deleted", nil)
}

func (h *DiscountHandler) GetDiscounts(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)

	filter := map[string]interface{}{}
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}
	if idStr := c.Query("id"); idStr != "" {
		if id, err := primitive.ObjectIDFromHex(idStr); err == nil {
			filter["id"] = id
		}
	}
	if pct, err := strconv.Atoi(c.Query("percentage")); err == nil && pct > 0 {
		filter["percentage"] = uint(pct)
	}

	discounts, total, err := h.discountUC.ListDiscounts(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		respondInternalError(c, err, "GetDiscounts")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, gin.H{
		"status": true, "data": discounts, "currentPage": page,
		"pageSize": pageSize, "totalPages": totalPages, "totalItems": total,
	})
}
