package usecase

import (
	"context"
	"errors"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartUsecase struct {
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
}

func NewCartUsecase(cr domain.CartRepository, pr domain.ProductRepository) *CartUsecase {
	return &CartUsecase{cartRepo: cr, productRepo: pr}
}

func (uc *CartUsecase) AddToCart(ctx context.Context, userID, productID primitive.ObjectID, quantity uint) (map[string]interface{}, error) {
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if product.Stock < quantity {
		return nil, errors.New("insufficient stock")
	}

	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Check if product already in cart
	existing, err := uc.cartRepo.FindProductInCart(ctx, cart.ID, productID)
	if err == nil && existing != nil {
		newQty := existing.Quantity + quantity
		newPrice := product.Price * newQty
		if err := uc.cartRepo.UpdateProductQuantity(ctx, existing.ID, newQty, newPrice); err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"product_name": product.Name,
			"quantity":     newQty,
			"price":        newPrice,
		}, nil
	}

	cp := &domain.CartProduct{
		CartID:       cart.ID,
		ProductID:    productID,
		ProductName:  product.Name,
		Quantity:     quantity,
		ProductPrice: product.Price * quantity,
	}
	if err := uc.cartRepo.AddProduct(ctx, cp); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"product_name": product.Name,
		"quantity":     quantity,
		"price":        cp.ProductPrice,
	}, nil
}

func (uc *CartUsecase) ViewCart(ctx context.Context, userID primitive.ObjectID) ([]domain.CartProduct, uint, error) {
	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, 0, errors.New("cart not found")
		}
		return nil, 0, err
	}

	items, err := uc.cartRepo.FindCartProducts(ctx, cart.ID)
	if err != nil {
		return nil, 0, err
	}

	var total uint
	for _, item := range items {
		total += item.ProductPrice
	}
	return items, total, nil
}
