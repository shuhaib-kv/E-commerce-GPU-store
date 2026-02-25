package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discount struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DiscountName       string             `bson:"discount_name" json:"discount_name"`
	DiscountPercentage uint               `bson:"discount_percentage" json:"discount_percentage"`
}

type DiscountRepository interface {
	Create(ctx context.Context, discount *Discount) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*Discount, error)
	FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]Discount, int64, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}
