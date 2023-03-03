package models

type RazorPay struct {
	UserID          uint
	RazorPaymentId  string
	RazorPayOrderID string
	Orderid         string
	Signature       string
	OrderId         string
	AmountPaid      string
}
