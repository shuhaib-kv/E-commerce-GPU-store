package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderRepo struct {
	orders *mongo.Collection
	items  *mongo.Collection
}

func NewOrderRepo(db *mongo.Database) domain.OrderRepository {
	return &orderRepo{
		orders: db.Collection("orders"),
		items:  db.Collection("order_items"),
	}
}

func (r *orderRepo) Create(ctx context.Context, order *domain.Order) error {
	res, err := r.orders.InsertOne(ctx, order)
	if err != nil {
		return err
	}
	order.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *orderRepo) CreateItem(ctx context.Context, item *domain.OrderItem) error {
	res, err := r.items.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *orderRepo) FindByOrderID(ctx context.Context, orderID string) (*domain.Order, error) {
	var order domain.Order
	err := r.orders.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Order, error) {
	cursor, err := r.orders.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepo) FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.Order, int64, error) {
	mongoFilter := bson.M{}
	for k, v := range filter {
		mongoFilter[k] = v
	}

	total, err := r.orders.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.orders.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *orderRepo) FindItemsByOrderID(ctx context.Context, orderID string) ([]domain.OrderItem, error) {
	cursor, err := r.items.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []domain.OrderItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *orderRepo) UpdateStatus(ctx context.Context, orderID string, status, paymentStatus bool) error {
	_, err := r.orders.UpdateOne(ctx, bson.M{"order_id": orderID}, bson.M{"$set": bson.M{
		"status":         status,
		"payment_status": paymentStatus,
	}})
	return err
}
