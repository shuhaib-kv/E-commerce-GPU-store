package usecase

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiscountUsecase struct {
	discountRepo domain.DiscountRepository
}

func NewDiscountUsecase(dr domain.DiscountRepository) *DiscountUsecase {
	return &DiscountUsecase{discountRepo: dr}
}

func (uc *DiscountUsecase) AddDiscount(ctx context.Context, name string, percentage uint) (*domain.Discount, error) {
	d := &domain.Discount{DiscountName: name, DiscountPercentage: percentage}
	if err := uc.discountRepo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (uc *DiscountUsecase) DeleteDiscount(ctx context.Context, id primitive.ObjectID) error {
	return uc.discountRepo.Delete(ctx, id)
}

func (uc *DiscountUsecase) ListDiscounts(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.Discount, int64, error) {
	return uc.discountRepo.FindAll(ctx, filter, page, pageSize)
}
