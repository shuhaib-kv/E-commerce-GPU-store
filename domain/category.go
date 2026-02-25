package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*Category, error)
	FindByName(ctx context.Context, name string) (*Category, error)
	FindAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, id primitive.ObjectID, name string) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
