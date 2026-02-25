package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type discountRepo struct {
	col *mongo.Collection
}

func NewDiscountRepo(db *mongo.Database) domain.DiscountRepository {
	return &discountRepo{col: db.Collection("discounts")}
}

func (r *discountRepo) Create(ctx context.Context, discount *domain.Discount) error {
	res, err := r.col.InsertOne(ctx, discount)
	if err != nil {
		return err
	}
	discount.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *discountRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Discount, error) {
	var discount domain.Discount
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&discount)
	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *discountRepo) FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.Discount, int64, error) {
	mongoFilter := bson.M{}
	if name, ok := filter["name"].(string); ok && name != "" {
		mongoFilter["discount_name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if id, ok := filter["id"].(primitive.ObjectID); ok {
		mongoFilter["_id"] = id
	}
	if pct, ok := filter["percentage"].(uint); ok && pct > 0 {
		mongoFilter["discount_percentage"] = pct
	}

	total, err := r.col.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.col.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var discounts []domain.Discount
	if err := cursor.All(ctx, &discounts); err != nil {
		return nil, 0, err
	}
	return discounts, total, nil
}

func (r *discountRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
