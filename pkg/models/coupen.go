package models

import "gorm.io/gorm"

type Coupon struct {
	gorm.Model
	CouponName       string`json:""`
	CouponCode       string`json:""`
	CouponPercentage int`json:""`
}
