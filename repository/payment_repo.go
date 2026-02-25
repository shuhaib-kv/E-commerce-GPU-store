package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type paymentRepo struct {
	col *mongo.Collection
}

func NewPaymentRepo(db *mongo.Database) domain.PaymentRepository {
	return &paymentRepo{col: db.Collection("payments")}
}

func (r *paymentRepo) Create(ctx context.Context, payment *domain.RazorPayment) error {
	res, err := r.col.InsertOne(ctx, payment)
	if err != nil {
		return err
	}
	payment.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *paymentRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.RazorPayment, error) {
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []domain.RazorPayment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepo) FindByOrderID(ctx context.Context, orderID string) (*domain.RazorPayment, error) {
	var payment domain.RazorPayment
	err := r.col.FindOne(ctx, bson.M{"razor_pay_order_id": orderID}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
