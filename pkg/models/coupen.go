package models

import "gorm.io/gorm"

type Coupon struct {
	gorm.Model
	CouponName       string
	CouponCode       string
	CouponPercentage uint
}
