package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Price          uint               `bson:"price" json:"price"`
	ModelNo        uint               `bson:"model_no" json:"model_no"`
	Image1         string             `bson:"image1" json:"image1"`
	Image2         string             `bson:"image2" json:"image2"`
	Image3         string             `bson:"image3" json:"image3"`
	Stock          uint               `bson:"stock" json:"stock"`
	CategoryID     primitive.ObjectID `bson:"category_id" json:"category_id"`
	Description    string             `bson:"description" json:"description"`
	Brand          string             `bson:"brand" json:"brand"`
	DiscountID     primitive.ObjectID `bson:"discount_id,omitempty" json:"discount_id"`
	Specifications bson.M             `bson:"specifications,omitempty" json:"specifications"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type ProductAttributeDefinition struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryID    primitive.ObjectID `bson:"category_id" json:"category_id"`
	AttributeKey  string             `bson:"attribute_key" json:"attribute_key"`
	DisplayName   string             `bson:"display_name" json:"display_name"`
	AttributeType string             `bson:"attribute_type" json:"attribute_type"`
	Required      bool               `bson:"required" json:"required"`
	Options       []string           `bson:"options,omitempty" json:"options"`
}

type ProductFilter struct {
	Name      string `form:"name"`
	Brand     string `form:"brand"`
	MinPrice  uint   `form:"minPrice"`
	MaxPrice  uint   `form:"maxPrice"`
	Category  string `form:"category"`
	PageSize  int    `form:"pageSize"`
	PageIndex int    `form:"pageIndex"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*Product, error)
	FindByModelNo(ctx context.Context, modelNo uint) (*Product, error)
	FindAll(ctx context.Context, filter bson.M, page, pageSize int) ([]Product, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	CreateAttributeDefinition(ctx context.Context, attr *ProductAttributeDefinition) error
	FindAttributesByCategory(ctx context.Context, categoryID primitive.ObjectID) ([]ProductAttributeDefinition, error)
	DeleteAttributeDefinition(ctx context.Context, id primitive.ObjectID) error
	FindAttributeByKeyAndCategory(ctx context.Context, categoryID primitive.ObjectID, key string) (*ProductAttributeDefinition, error)
}
