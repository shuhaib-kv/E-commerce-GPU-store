package usecase

import (
	"context"
	"errors"
	"ga/domain"
)

type AdminUsecase struct {
	adminRepo domain.AdminRepository
}

func NewAdminUsecase(ar domain.AdminRepository) *AdminUsecase {
	return &AdminUsecase{adminRepo: ar}
}

func (uc *AdminUsecase) Signup(ctx context.Context, admin *domain.Admin) error {
	if existing, _ := uc.adminRepo.FindByEmail(ctx, admin.Email); existing != nil {
		return errors.New("email already exists")
	}
	if err := admin.HashPassword(admin.Password); err != nil {
		return err
	}
	return uc.adminRepo.Create(ctx, admin)
}

func (uc *AdminUsecase) Login(ctx context.Context, email, password string) (*domain.Admin, error) {
	admin, err := uc.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !admin.CheckPassword(password) {
		return nil, errors.New("invalid email or password")
	}
	return admin, nil
}
