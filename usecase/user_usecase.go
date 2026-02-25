package usecase

import (
	"context"
	"errors"
	"ga/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserUsecase struct {
	userRepo   domain.UserRepository
	walletRepo domain.WalletRepository
	cartRepo   domain.CartRepository
}

func NewUserUsecase(ur domain.UserRepository, wr domain.WalletRepository, cr domain.CartRepository) *UserUsecase {
	return &UserUsecase{userRepo: ur, walletRepo: wr, cartRepo: cr}
}

func (uc *UserUsecase) Signup(ctx context.Context, user *domain.User) error {
	if existing, _ := uc.userRepo.FindByEmail(ctx, user.Email); existing != nil {
		return errors.New("email already exists")
	}
	if existing, _ := uc.userRepo.FindByUsername(ctx, user.UserName); existing != nil {
		return errors.New("username already exists")
	}

	if err := user.HashPassword(user.Password); err != nil {
		return err
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Create wallet
	uc.walletRepo.Create(ctx, &domain.Wallet{UserID: user.ID, Balance: 0})
	// Create cart
	uc.cartRepo.Create(ctx, &domain.Cart{UserID: user.ID})

	return nil
}

func (uc *UserUsecase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if user.BlockStatus {
		return nil, errors.New("user is blocked")
	}
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

func (uc *UserUsecase) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return uc.userRepo.FindByPhone(ctx, phone)
}

func (uc *UserUsecase) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

func (uc *UserUsecase) AddAddress(ctx context.Context, address *domain.Address) error {
	return uc.userRepo.CreateAddress(ctx, address)
}

func (uc *UserUsecase) GetAddress(ctx context.Context, userID primitive.ObjectID) (*domain.Address, error) {
	return uc.userRepo.FindAddressByUserID(ctx, userID)
}

func (uc *UserUsecase) EditAddress(ctx context.Context, userID primitive.ObjectID, update map[string]interface{}) error {
	return uc.userRepo.UpdateAddress(ctx, userID, update)
}

// Admin user management
func (uc *UserUsecase) ListUsers(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.User, int64, error) {
	return uc.userRepo.FindAll(ctx, filter, page, pageSize)
}

func (uc *UserUsecase) BlockUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("user not found")
		}
		return err
	}
	return uc.userRepo.Update(ctx, id, map[string]interface{}{"block_status": true})
}

func (uc *UserUsecase) UnblockUser(ctx context.Context, id primitive.ObjectID) error {
	return uc.userRepo.Update(ctx, id, map[string]interface{}{"block_status": false})
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return uc.userRepo.Delete(ctx, id)
}
