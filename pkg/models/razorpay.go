package models

type RazorPay struct {
	UserID          int
	RazorPaymentId  string
	RazorPayOrderID string
	Signature       string
	OrderId         string
	AmountPaid      string
}
