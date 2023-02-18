package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddCoupon(c *gin.Context) {
	// Get Info from the req body

	couponName := c.PostForm("couponName")
	couponCode := c.PostForm("couponCode")
	Percentage := c.PostForm("couponPercentage")
	couponPercentage, _ := strconv.Atoi(Percentage)

	coupon := models.Coupon{CouponName: couponName, CouponCode: couponCode, CouponPercentage: couponPercentage}

	var checkCoup []models.Coupon
	database.Db.Find(&checkCoup)

	// Checking username existence
	for _, i := range checkCoup {
		if i.CouponName == couponName {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Coupon Name Already Exist",
			})
			return
		}
	}

	result := database.Db.Create(&coupon)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error Creating Coupon",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Coupon Crearted",
	})

}

func DeleteCoupon(c *gin.Context) {
	var coupon models.Coupon
	couponName := c.Query("couponName")
	database.Db.Where("coupon_name = ?", couponName).Delete(&coupon)
	//database.DB.Raw("DELETE FROM coupons WHERE coupon_name=?", couponName).Scan(&coupon)
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Deleted succesfully",
	})

}

func ListCoupons(c *gin.Context) {
	var coupons []models.Coupon
	result := database.Db.Find(&coupons)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"ststus":  false,
			"message": "No coupon found",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": coupons,
	})

}
