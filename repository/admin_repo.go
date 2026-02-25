package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type adminRepo struct {
	col *mongo.Collection
}

func NewAdminRepo(db *mongo.Database) domain.AdminRepository {
	return &adminRepo{col: db.Collection("admins")}
}

func (r *adminRepo) Create(ctx context.Context, admin *domain.Admin) error {
	res, err := r.col.InsertOne(ctx, admin)
	if err != nil {
		return err
	}
	admin.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *adminRepo) FindByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	var admin domain.Admin
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
