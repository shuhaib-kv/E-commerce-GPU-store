package models

import "gorm.io/gorm"

type RazorPay struct {
	gorm.Model
	UserID          uint
	RazorPaymentId  string
	RazorPayOrderID string
	Signature       string
	AmountPaid      uint
}
