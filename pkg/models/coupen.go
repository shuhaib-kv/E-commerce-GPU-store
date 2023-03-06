package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	CouponName       string
	CouponCode       string
	CouponPercentage uint
	ExpiryDate       time.Time
}
