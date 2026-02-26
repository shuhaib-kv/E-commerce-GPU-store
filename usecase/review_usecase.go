package usecase

import (
	"context"
	"errors"
	"ga/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewUsecase struct {
	reviewRepo domain.ReviewRepository
}

func NewReviewUsecase(rr domain.ReviewRepository) *ReviewUsecase {
	return &ReviewUsecase{reviewRepo: rr}
}

func (uc *ReviewUsecase) AddReview(ctx context.Context, review *domain.Review) error {
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	existing, err := uc.reviewRepo.FindByUserAndProduct(ctx, review.UserID, review.ProductID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("you have already reviewed this product")
	}
	review.CreatedAt = time.Now()
	return uc.reviewRepo.Create(ctx, review)
}

func (uc *ReviewUsecase) GetProductReviews(ctx context.Context, productID primitive.ObjectID) ([]domain.Review, float64, error) {
	reviews, err := uc.reviewRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, 0, err
	}
	avg, err := uc.reviewRepo.AverageRating(ctx, productID)
	if err != nil {
		return nil, 0, err
	}
	return reviews, avg, nil
}
