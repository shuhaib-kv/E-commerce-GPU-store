package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewRepo struct {
	reviews *mongo.Collection
}

func NewReviewRepo(db *mongo.Database) domain.ReviewRepository {
	return &reviewRepo{reviews: db.Collection("reviews")}
}

func (r *reviewRepo) Create(ctx context.Context, review *domain.Review) error {
	res, err := r.reviews.InsertOne(ctx, review)
	if err != nil {
		return err
	}
	review.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *reviewRepo) FindByProductID(ctx context.Context, productID primitive.ObjectID) ([]domain.Review, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.reviews.Find(ctx, bson.M{"product_id": productID}, opts)
	if err != nil {
		return nil, err
	}
	var reviews []domain.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepo) FindByUserAndProduct(ctx context.Context, userID, productID primitive.ObjectID) (*domain.Review, error) {
	var review domain.Review
	err := r.reviews.FindOne(ctx, bson.M{"user_id": userID, "product_id": productID}).Decode(&review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepo) AverageRating(ctx context.Context, productID primitive.ObjectID) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"product_id": productID}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "avg": bson.M{"$avg": "$rating"}}}},
	}
	cursor, err := r.reviews.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, nil
	}
	avg, _ := results[0]["avg"].(float64)
	return avg, nil
}
