package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CouponName       string             `bson:"coupon_name" json:"coupon_name"`
	CouponCode       string             `bson:"coupon_code" json:"coupon_code"`
	CouponPercentage uint               `bson:"coupon_percentage" json:"coupon_percentage"`
	ExpiryDate       time.Time          `bson:"expiry_date" json:"expiry_date"`
}

type CouponRepository interface {
	Create(ctx context.Context, coupon *Coupon) error
	FindByName(ctx context.Context, name string) (*Coupon, error)
	FindByCode(ctx context.Context, code string) (*Coupon, error)
	FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]Coupon, int64, error)
	Delete(ctx context.Context, name string) error
}
