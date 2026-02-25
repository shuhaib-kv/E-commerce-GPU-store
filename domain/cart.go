package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}

type CartProduct struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CartID       primitive.ObjectID `bson:"cart_id" json:"cart_id"`
	ProductID    primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductName  string             `bson:"product_name" json:"product_name"`
	Quantity     uint               `bson:"quantity" json:"quantity"`
	ProductPrice uint               `bson:"product_price" json:"product_price"`
}

type CartRepository interface {
	Create(ctx context.Context, cart *Cart) error
	FindByUserID(ctx context.Context, userID primitive.ObjectID) (*Cart, error)
	AddProduct(ctx context.Context, cp *CartProduct) error
	FindProductInCart(ctx context.Context, cartID, productID primitive.ObjectID) (*CartProduct, error)
	UpdateProductQuantity(ctx context.Context, id primitive.ObjectID, quantity uint, price uint) error
	FindCartProducts(ctx context.Context, cartID primitive.ObjectID) ([]CartProduct, error)
	DeleteCartProducts(ctx context.Context, cartID primitive.ObjectID) error
}
