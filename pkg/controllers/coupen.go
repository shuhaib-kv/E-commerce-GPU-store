package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddCoupon(c *gin.Context) {
	// Get Info from the req body
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

	// Checking username existence
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
	var coupon models.Coupon
	couponName := c.Query("couponName")
	database.Db.Where("coupon_name = ?", couponName).Delete(&coupon)
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Deleted succesfully",
	})

}

func ListCoupons(c *gin.Context) {
	var coupons []models.Coupon
	result := database.Db.Find(&coupons)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"ststus":  false,
			"message": "No coupon found",
			"error":   result.Error,
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"status":  true,
		"message": "coupen found",
		"data":    coupons,
	})

}
