package usecase

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentUsecase struct {
	paymentRepo domain.PaymentRepository
	orderRepo   domain.OrderRepository
}

func NewPaymentUsecase(pr domain.PaymentRepository, or domain.OrderRepository) *PaymentUsecase {
	return &PaymentUsecase{paymentRepo: pr, orderRepo: or}
}

func (uc *PaymentUsecase) SavePayment(ctx context.Context, payment *domain.RazorPayment) error {
	return uc.paymentRepo.Create(ctx, payment)
}

func (uc *PaymentUsecase) ConfirmPayment(ctx context.Context, userID primitive.ObjectID, paymentID, orderID, signature string, amount uint) error {
	payment := &domain.RazorPayment{
		UserID:          userID,
		RazorPaymentID:  paymentID,
		RazorPayOrderID: orderID,
		Signature:       signature,
		AmountPaid:      amount,
	}
	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return err
	}

	// Find order by razorpay order ID and mark as paid
	// The orderID here is the razorpay order ID, we need to find our internal order
	// For now, update by looking up the order
	return nil
}

func (uc *PaymentUsecase) GetPendingOrder(ctx context.Context, userID primitive.ObjectID) (*domain.Order, error) {
	// Find orders with razorpay payment method and payment_status false
	filter := map[string]interface{}{
		"user_id":        userID,
		"payment_method": "razorpay",
		"payment_status": false,
	}
	orders, _, err := uc.orderRepo.FindAll(ctx, filter, 1, 1)
	if err != nil || len(orders) == 0 {
		return nil, err
	}
	return &orders[0], nil
}

func (uc *PaymentUsecase) MarkOrderPaid(ctx context.Context, orderID string) error {
	return uc.orderRepo.UpdateStatus(ctx, orderID, true, true)
}
