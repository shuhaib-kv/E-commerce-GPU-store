package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryRepo struct {
	col *mongo.Collection
}

func NewCategoryRepo(db *mongo.Database) domain.CategoryRepository {
	return &categoryRepo{col: db.Collection("categories")}
}

func (r *categoryRepo) Create(ctx context.Context, category *domain.Category) error {
	res, err := r.col.InsertOne(ctx, category)
	if err != nil {
		return err
	}
	category.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *categoryRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	var cat domain.Category
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepo) FindByName(ctx context.Context, name string) (*domain.Category, error) {
	var cat domain.Category
	err := r.col.FindOne(ctx, bson.M{"name": name}).Decode(&cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepo) FindAll(ctx context.Context) ([]domain.Category, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) Update(ctx context.Context, id primitive.ObjectID, name string) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name}})
	return err
}

func (r *categoryRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
