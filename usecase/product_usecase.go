package usecase

import (
	"context"
	"errors"
	"ga/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductUsecase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
	discountRepo domain.DiscountRepository
}

func NewProductUsecase(pr domain.ProductRepository, cr domain.CategoryRepository, dr domain.DiscountRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: pr, categoryRepo: cr, discountRepo: dr}
}

func (uc *ProductUsecase) AddProduct(ctx context.Context, product *domain.Product) error {
	if existing, _ := uc.productRepo.FindByModelNo(ctx, product.ModelNo); existing != nil {
		return errors.New("product with this model number already exists")
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	return uc.productRepo.Create(ctx, product)
}

func (uc *ProductUsecase) EditProduct(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("product not found")
		}
		return err
	}
	update["updated_at"] = time.Now()
	return uc.productRepo.Update(ctx, id, update)
}

func (uc *ProductUsecase) ViewProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	mongoFilter := bson.M{}
	if filter.Name != "" {
		mongoFilter["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}
	if filter.Brand != "" {
		mongoFilter["brand"] = filter.Brand
	}
	if filter.MinPrice > 0 {
		mongoFilter["price"] = bson.M{"$gte": filter.MinPrice}
	}
	if filter.MaxPrice > 0 {
		if existing, ok := mongoFilter["price"].(bson.M); ok {
			existing["$lte"] = filter.MaxPrice
		} else {
			mongoFilter["price"] = bson.M{"$lte": filter.MaxPrice}
		}
	}
	if filter.Category != "" {
		cat, err := uc.categoryRepo.FindByName(ctx, filter.Category)
		if err == nil {
			mongoFilter["category_id"] = cat.ID
		}
	}

	page := filter.PageIndex
	if page == 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	return uc.productRepo.FindAll(ctx, mongoFilter, page, pageSize)
}

func (uc *ProductUsecase) ViewProductsUser(ctx context.Context, filter bson.M, page, pageSize int) ([]map[string]interface{}, int64, error) {
	products, total, err := uc.productRepo.FindAll(ctx, filter, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var response []map[string]interface{}
	for _, product := range products {
		price := product.Price
		discountAmount := uint(0)
		if !product.DiscountID.IsZero() {
			discount, err := uc.discountRepo.FindByID(ctx, product.DiscountID)
			if err == nil {
				discountAmount = uint(float64(price) * float64(discount.DiscountPercentage) / 100.0)
				price = price - discountAmount
			}
		}
		response = append(response, map[string]interface{}{
			"id":              product.ID,
			"name":            product.Name,
			"price":           price,
			"image1":          product.Image1,
			"image2":          product.Image2,
			"image3":          product.Image3,
			"brand":           product.Brand,
			"specifications":  product.Specifications,
			"discount_amount": discountAmount,
		})
	}
	return response, total, nil
}

func (uc *ProductUsecase) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	_, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("product not found")
	}
	return uc.productRepo.Delete(ctx, id)
}

func (uc *ProductUsecase) GetProduct(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	return uc.productRepo.FindByID(ctx, id)
}

// Attribute definitions
func (uc *ProductUsecase) AddAttributeDefinition(ctx context.Context, attr *domain.ProductAttributeDefinition) error {
	// Verify category exists
	_, err := uc.categoryRepo.FindByID(ctx, attr.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}
	// Check duplicate
	if existing, _ := uc.productRepo.FindAttributeByKeyAndCategory(ctx, attr.CategoryID, attr.AttributeKey); existing != nil {
		return errors.New("attribute key already exists for this category")
	}
	return uc.productRepo.CreateAttributeDefinition(ctx, attr)
}

func (uc *ProductUsecase) ListAttributeDefinitions(ctx context.Context, categoryID primitive.ObjectID) ([]domain.ProductAttributeDefinition, error) {
	return uc.productRepo.FindAttributesByCategory(ctx, categoryID)
}

func (uc *ProductUsecase) DeleteAttributeDefinition(ctx context.Context, id primitive.ObjectID) error {
	return uc.productRepo.DeleteAttributeDefinition(ctx, id)
}
