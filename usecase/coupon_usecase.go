package usecase

import (
	"context"
	"errors"
	"ga/domain"
	"time"

	"github.com/google/uuid"
)

type CouponUsecase struct {
	couponRepo domain.CouponRepository
}

func NewCouponUsecase(cr domain.CouponRepository) *CouponUsecase {
	return &CouponUsecase{couponRepo: cr}
}

func (uc *CouponUsecase) AddCoupon(ctx context.Context, name string, percentage, expiresInDays uint) (*domain.Coupon, error) {
	if existing, _ := uc.couponRepo.FindByName(ctx, name); existing != nil {
		return nil, errors.New("coupon name already exists")
	}

	coupon := &domain.Coupon{
		CouponName:       name,
		CouponCode:       "STORE-" + uuid.New().String(),
		CouponPercentage: percentage,
		ExpiryDate:       time.Now().AddDate(0, 0, int(expiresInDays)),
	}

	if err := uc.couponRepo.Create(ctx, coupon); err != nil {
		return nil, err
	}
	return coupon, nil
}

func (uc *CouponUsecase) DeleteCoupon(ctx context.Context, name string) error {
	return uc.couponRepo.Delete(ctx, name)
}

func (uc *CouponUsecase) ListCoupons(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.Coupon, int64, error) {
	return uc.couponRepo.FindAll(ctx, filter, page, pageSize)
}
