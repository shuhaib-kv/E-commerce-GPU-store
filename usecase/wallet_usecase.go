package usecase

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletUsecase struct {
	walletRepo domain.WalletRepository
}

func NewWalletUsecase(wr domain.WalletRepository) *WalletUsecase {
	return &WalletUsecase{walletRepo: wr}
}

func (uc *WalletUsecase) GetWalletInfo(ctx context.Context, userID primitive.ObjectID) (*domain.Wallet, []domain.WalletHistory, error) {
	wallet, err := uc.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	history, err := uc.walletRepo.FindHistoryByUserID(ctx, userID)
	if err != nil {
		return wallet, nil, nil
	}
	return wallet, history, nil
}
