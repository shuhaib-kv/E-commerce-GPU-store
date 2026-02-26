package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	UserName  string             `bson:"user_name" json:"user_name"`
	Rating    int                `bson:"rating" json:"rating"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	FindByProductID(ctx context.Context, productID primitive.ObjectID) ([]Review, error)
	FindByUserAndProduct(ctx context.Context, userID, productID primitive.ObjectID) (*Review, error)
	AverageRating(ctx context.Context, productID primitive.ObjectID) (float64, error)
}
