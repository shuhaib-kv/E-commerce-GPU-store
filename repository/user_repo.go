package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	users     *mongo.Collection
	addresses *mongo.Collection
}

func NewUserRepo(db *mongo.Database) domain.UserRepository {
	return &userRepo{
		users:     db.Collection("users"),
		addresses: db.Collection("addresses"),
	}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	res, err := r.users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.users.FindOne(ctx, bson.M{"user_name": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	err := r.users.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.User, int64, error) {
	mongoFilter := bson.M{}
	if name, ok := filter["name"].(string); ok && name != "" {
		mongoFilter["first_name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if id, ok := filter["id"].(primitive.ObjectID); ok {
		mongoFilter["_id"] = id
	}

	total, err := r.users.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.users.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error {
	_, err := r.users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *userRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.users.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *userRepo) CreateAddress(ctx context.Context, address *domain.Address) error {
	res, err := r.addresses.InsertOne(ctx, address)
	if err != nil {
		return err
	}
	address.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *userRepo) FindAddressByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Address, error) {
	var addr domain.Address
	err := r.addresses.FindOne(ctx, bson.M{"user_id": userID}).Decode(&addr)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (r *userRepo) UpdateAddress(ctx context.Context, userID primitive.ObjectID, update map[string]interface{}) error {
	_, err := r.addresses.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": update})
	return err
}
