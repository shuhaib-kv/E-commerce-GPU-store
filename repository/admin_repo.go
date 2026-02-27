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

func (r *adminRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Admin, error) {
	var admin domain.Admin
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepo) FindAll(ctx context.Context) ([]*domain.Admin, error) {
	cursor, err := r.col.Find(ctx, bson.M{"role": bson.M{"$ne": "superadmin"}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var admins []*domain.Admin
	if err := cursor.All(ctx, &admins); err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *adminRepo) Update(ctx context.Context, admin *domain.Admin) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": admin.ID}, bson.M{"$set": admin})
	return err
}
