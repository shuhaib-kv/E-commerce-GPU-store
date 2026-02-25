package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cartRepo struct {
	carts    *mongo.Collection
	products *mongo.Collection
}

func NewCartRepo(db *mongo.Database) domain.CartRepository {
	return &cartRepo{
		carts:    db.Collection("carts"),
		products: db.Collection("cart_products"),
	}
}

func (r *cartRepo) Create(ctx context.Context, cart *domain.Cart) error {
	res, err := r.carts.InsertOne(ctx, cart)
	if err != nil {
		return err
	}
	cart.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *cartRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.carts.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepo) AddProduct(ctx context.Context, cp *domain.CartProduct) error {
	res, err := r.products.InsertOne(ctx, cp)
	if err != nil {
		return err
	}
	cp.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *cartRepo) FindProductInCart(ctx context.Context, cartID, productID primitive.ObjectID) (*domain.CartProduct, error) {
	var cp domain.CartProduct
	err := r.products.FindOne(ctx, bson.M{"cart_id": cartID, "product_id": productID}).Decode(&cp)
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *cartRepo) UpdateProductQuantity(ctx context.Context, id primitive.ObjectID, quantity uint, price uint) error {
	_, err := r.products.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"quantity": quantity, "product_price": price}})
	return err
}

func (r *cartRepo) FindCartProducts(ctx context.Context, cartID primitive.ObjectID) ([]domain.CartProduct, error) {
	cursor, err := r.products.Find(ctx, bson.M{"cart_id": cartID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []domain.CartProduct
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartRepo) DeleteCartProducts(ctx context.Context, cartID primitive.ObjectID) error {
	_, err := r.products.DeleteMany(ctx, bson.M{"cart_id": cartID})
	return err
}
