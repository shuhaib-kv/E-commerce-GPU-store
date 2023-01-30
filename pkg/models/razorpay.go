package models

type RazorPay struct {
	UserID          int
	RazorPaymentId  string`json:""`
	RazorPayOrderID string`json:""`
	Signature       string`json:""`
	OrderId         string`json:""`
	AmountPaid      string`json:""`
}
