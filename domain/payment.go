package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RazorPayment struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	RazorPaymentID  string             `bson:"razor_payment_id" json:"razor_payment_id"`
	RazorPayOrderID string             `bson:"razor_pay_order_id" json:"razor_pay_order_id"`
	Signature       string             `bson:"signature" json:"signature"`
	AmountPaid      uint               `bson:"amount_paid" json:"amount_paid"`
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *RazorPayment) error
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]RazorPayment, error)
	FindByOrderID(ctx context.Context, orderID string) (*RazorPayment, error)
}
