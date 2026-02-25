package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepo struct {
	products *mongo.Collection
	attrDefs *mongo.Collection
}

func NewProductRepo(db *mongo.Database) domain.ProductRepository {
	return &productRepo{
		products: db.Collection("products"),
		attrDefs: db.Collection("product_attribute_definitions"),
	}
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	res, err := r.products.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	product.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *productRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	err := r.products.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) FindByModelNo(ctx context.Context, modelNo uint) (*domain.Product, error) {
	var product domain.Product
	err := r.products.FindOne(ctx, bson.M{"model_no": modelNo}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) FindAll(ctx context.Context, filter bson.M, page, pageSize int) ([]domain.Product, int64, error) {
	total, err := r.products.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.products.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *productRepo) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.products.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *productRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.products.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *productRepo) CreateAttributeDefinition(ctx context.Context, attr *domain.ProductAttributeDefinition) error {
	res, err := r.attrDefs.InsertOne(ctx, attr)
	if err != nil {
		return err
	}
	attr.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *productRepo) FindAttributesByCategory(ctx context.Context, categoryID primitive.ObjectID) ([]domain.ProductAttributeDefinition, error) {
	cursor, err := r.attrDefs.Find(ctx, bson.M{"category_id": categoryID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attrs []domain.ProductAttributeDefinition
	if err := cursor.All(ctx, &attrs); err != nil {
		return nil, err
	}
	return attrs, nil
}

func (r *productRepo) DeleteAttributeDefinition(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.attrDefs.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *productRepo) FindAttributeByKeyAndCategory(ctx context.Context, categoryID primitive.ObjectID, key string) (*domain.ProductAttributeDefinition, error) {
	var attr domain.ProductAttributeDefinition
	err := r.attrDefs.FindOne(ctx, bson.M{"category_id": categoryID, "attribute_key": key}).Decode(&attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}
