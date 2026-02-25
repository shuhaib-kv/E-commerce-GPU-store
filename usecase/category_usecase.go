package usecase

import (
	"context"
	"errors"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryUsecase struct {
	categoryRepo domain.CategoryRepository
	productRepo  domain.ProductRepository
}

func NewCategoryUsecase(cr domain.CategoryRepository, pr domain.ProductRepository) *CategoryUsecase {
	return &CategoryUsecase{categoryRepo: cr, productRepo: pr}
}

func (uc *CategoryUsecase) AddCategory(ctx context.Context, name string) (*domain.Category, error) {
	if existing, _ := uc.categoryRepo.FindByName(ctx, name); existing != nil {
		return nil, errors.New("category already exists")
	}
	cat := &domain.Category{Name: name}
	if err := uc.categoryRepo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (uc *CategoryUsecase) ViewAll(ctx context.Context) ([]domain.Category, error) {
	return uc.categoryRepo.FindAll(ctx)
}

func (uc *CategoryUsecase) EditCategory(ctx context.Context, id primitive.ObjectID, name string) error {
	return uc.categoryRepo.Update(ctx, id, name)
}

func (uc *CategoryUsecase) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	return uc.categoryRepo.Delete(ctx, id)
}

func (uc *CategoryUsecase) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	return uc.categoryRepo.FindByID(ctx, id)
}
