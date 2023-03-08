package controllers

import (
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddCoupon(c *gin.Context) {
	var body struct {
		CouponName       string `json:"couponname"`
		CouponCode       string `json:"couponcode"`
		CouponPercentage uint   `json:"couponpercentage"`
		Expiresat        uint   `json:"expiresat"`
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
	coupon := models.Coupon{CouponName: body.CouponName,
		CouponCode:       body.CouponCode,
		CouponPercentage: body.CouponPercentage,
		ExpiryDate:       time.Now().AddDate(0, 0, int(body.Expiresat)),
	}

	var checkCoup []models.Coupon
	database.Db.Find(&checkCoup)

	for _, i := range checkCoup {
		if i.CouponName == body.CouponName {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"error":   "Coupon Name Already Exist",
				"message": "Duplicate Coupen Name",
			})
			return
		}
	}

	result := database.Db.Create(&coupon)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"error":   result.Error,
			"message": "Error Creating Coupon",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Coupon Crearted",
		"data":    coupon,
	})

}

func DeleteCoupon(c *gin.Context) {
	var body struct {
		CouponName string `json:"couponname"`
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

	result := database.Db.Delete(&models.Coupon{}, "coupon_name = ?", body.CouponName)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to delete coupon",
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Coupon not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Coupon deleted successfully",
	})
}
func ListCoupons(c *gin.Context) {
	var coupons []models.Coupon

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// Retrieve query parameters for search filter
	couponCode := c.Query("couponCode")
	couponName := c.Query("couponName")
	couponPercentage, err := strconv.ParseUint(c.Query("couponPercentage"), 10, 32)
	if err != nil {
		couponPercentage = 0
	}

	// Build query with search filter
	query := database.Db.Offset(offset).Limit(pageSize)
	if couponCode != "" {
		query = query.Where("coupon_code LIKE ?", fmt.Sprintf("%%%s%%", couponCode))
	}
	if couponName != "" {
		query = query.Where("coupon_name LIKE ?", fmt.Sprintf("%%%s%%", couponName))
	}
	if couponPercentage > 0 {
		query = query.Where("coupon_percentage = ?", couponPercentage)
	}

	result := query.Find(&coupons)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Error retrieving coupons",
			"error":   result.Error.Error(),
		})
		return
	}

	var totalCoupons int64
	database.Db.Model(&models.Coupon{}).Count(&totalCoupons)
	totalPages := int(math.Ceil(float64(totalCoupons) / float64(pageSize)))

	var couponsResponse []map[string]interface{}
	for _, coupon := range coupons {
		couponResponse := map[string]interface{}{
			"CouponName":       coupon.CouponName,
			"CouponCode":       coupon.CouponCode,
			"CouponPercentage": coupon.CouponPercentage,
			"ExpiryDate":       coupon.ExpiryDate,
		}
		couponsResponse = append(couponsResponse, couponResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"message":     "Coupons retrieved successfully",
		"data":        couponsResponse,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
		"totalItems":  totalCoupons,
	})
}
