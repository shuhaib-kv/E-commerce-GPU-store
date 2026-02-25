package handler

import (
	"ga/usecase"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	couponUC *usecase.CouponUsecase
}

func NewCouponHandler(cuc *usecase.CouponUsecase) *CouponHandler {
	return &CouponHandler{couponUC: cuc}
}

func (h *CouponHandler) AddCoupon(c *gin.Context) {
	var body struct {
		CouponName       string `json:"couponname" binding:"required"`
		CouponPercentage uint   `json:"couponpercentage" binding:"required"`
		ExpiresAt        uint   `json:"expiresat" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	coupon, err := h.couponUC.AddCoupon(c.Request.Context(), body.CouponName, body.CouponPercentage, body.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": true, "message": "Coupon created", "data": coupon})
}

func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	var body struct {
		CouponName string `json:"couponname" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}
	if err := h.couponUC.DeleteCoupon(c.Request.Context(), body.CouponName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "Coupon not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Coupon deleted"})
}

func (h *CouponHandler) ListCoupons(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)

	filter := map[string]interface{}{}
	if code := c.Query("couponCode"); code != "" {
		filter["coupon_code"] = code
	}
	if name := c.Query("couponName"); name != "" {
		filter["coupon_name"] = name
	}
	if pct, err := strconv.ParseUint(c.Query("couponPercentage"), 10, 32); err == nil && pct > 0 {
		filter["coupon_percentage"] = pct
	}

	coupons, total, err := h.couponUC.ListCoupons(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, gin.H{
		"status": true, "data": coupons, "currentPage": page,
		"pageSize": pageSize, "totalPages": totalPages, "totalItems": total,
	})
}
