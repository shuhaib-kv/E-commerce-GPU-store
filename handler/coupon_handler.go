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
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	coupon, err := h.couponUC.AddCoupon(c.Request.Context(), body.CouponName, body.CouponPercentage, body.ExpiresAt)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, "Coupon created", coupon)
}

func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	var body struct {
		CouponName string `json:"couponname" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.couponUC.DeleteCoupon(c.Request.Context(), body.CouponName); err != nil {
		respondError(c, http.StatusNotFound, "Coupon not found")
		return
	}
	respondSuccess(c, http.StatusOK, "Coupon deleted", nil)
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
		respondInternalError(c, err, "ListCoupons")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, gin.H{
		"status": true, "data": coupons, "currentPage": page,
		"pageSize": pageSize, "totalPages": totalPages, "totalItems": total,
	})
}
