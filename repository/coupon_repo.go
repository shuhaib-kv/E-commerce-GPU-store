package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type couponRepo struct {
	col *mongo.Collection
}

func NewCouponRepo(db *mongo.Database) domain.CouponRepository {
	return &couponRepo{col: db.Collection("coupons")}
}

func (r *couponRepo) Create(ctx context.Context, coupon *domain.Coupon) error {
	res, err := r.col.InsertOne(ctx, coupon)
	if err != nil {
		return err
	}
	coupon.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *couponRepo) FindByName(ctx context.Context, name string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := r.col.FindOne(ctx, bson.M{"coupon_name": name}).Decode(&coupon)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *couponRepo) FindByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := r.col.FindOne(ctx, bson.M{"coupon_code": code}).Decode(&coupon)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *couponRepo) FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.Coupon, int64, error) {
	mongoFilter := bson.M{}
	if code, ok := filter["coupon_code"].(string); ok && code != "" {
		mongoFilter["coupon_code"] = bson.M{"$regex": code, "$options": "i"}
	}
	if name, ok := filter["coupon_name"].(string); ok && name != "" {
		mongoFilter["coupon_name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if pct, ok := filter["coupon_percentage"].(uint64); ok && pct > 0 {
		mongoFilter["coupon_percentage"] = pct
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

	var coupons []domain.Coupon
	if err := cursor.All(ctx, &coupons); err != nil {
		return nil, 0, err
	}
	return coupons, total, nil
}

func (r *couponRepo) Delete(ctx context.Context, name string) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"coupon_name": name})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
