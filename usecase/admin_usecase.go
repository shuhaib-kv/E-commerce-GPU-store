package usecase

import (
	"context"
	"errors"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if admin.Role == "" {
		admin.Role = "admin"
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
	if admin.BlockStatus {
		return nil, errors.New("your account has been blocked")
	}
	if !admin.CheckPassword(password) {
		return nil, errors.New("invalid email or password")
	}
	return admin, nil
}

func (uc *AdminUsecase) ListAdmins(ctx context.Context) ([]*domain.Admin, error) {
	return uc.adminRepo.FindAll(ctx)
}

func (uc *AdminUsecase) GetAdmin(ctx context.Context, id string) (*domain.Admin, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid admin ID")
	}
	return uc.adminRepo.FindByID(ctx, objID)
}

func (uc *AdminUsecase) UpdateAdmin(ctx context.Context, admin *domain.Admin) error {
	return uc.adminRepo.Update(ctx, admin)
}

func (uc *AdminUsecase) BlockAdmin(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid admin ID")
	}
	admin, err := uc.adminRepo.FindByID(ctx, objID)
	if err != nil {
		return errors.New("admin not found")
	}
	if admin.Role == "superadmin" {
		return errors.New("cannot block a super admin")
	}
	admin.BlockStatus = true
	return uc.adminRepo.Update(ctx, admin)
}

func (uc *AdminUsecase) UnblockAdmin(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid admin ID")
	}
	admin, err := uc.adminRepo.FindByID(ctx, objID)
	if err != nil {
		return errors.New("admin not found")
	}
	admin.BlockStatus = false
	return uc.adminRepo.Update(ctx, admin)
}
