package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID               primitive.ObjectID `bson:"user_id" json:"user_id"`
	AddressID            primitive.ObjectID `bson:"address_id" json:"address_id"`
	OrderID              string             `bson:"order_id" json:"order_id"` // UUID string
	PaymentMethod        string             `bson:"payment_method" json:"payment_method"`
	TotalAmount          uint               `bson:"total_amount" json:"total_amount"`
	Status               bool               `bson:"status" json:"status"`
	PaymentStatus        bool               `bson:"payment_status" json:"payment_status"`
	ExpectedDeliveryDate time.Time          `bson:"expected_delivery_date" json:"expected_delivery_date"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
}

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID     string             `bson:"order_id" json:"order_id"`
	ProductID   primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductName string             `bson:"product_name" json:"product_name"`
	Quantity    uint               `bson:"quantity" json:"quantity"`
	Price       uint               `bson:"price" json:"price"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	CreateItem(ctx context.Context, item *OrderItem) error
	FindByOrderID(ctx context.Context, orderID string) (*Order, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]Order, error)
	FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]Order, int64, error)
	FindItemsByOrderID(ctx context.Context, orderID string) ([]OrderItem, error)
	UpdateStatus(ctx context.Context, orderID string, status, paymentStatus bool) error
}
