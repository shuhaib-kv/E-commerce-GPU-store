package repository

import (
	"context"
	"ga/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type walletRepo struct {
	wallets *mongo.Collection
	history *mongo.Collection
}

func NewWalletRepo(db *mongo.Database) domain.WalletRepository {
	return &walletRepo{
		wallets: db.Collection("wallets"),
		history: db.Collection("wallet_history"),
	}
}

func (r *walletRepo) Create(ctx context.Context, wallet *domain.Wallet) error {
	res, err := r.wallets.InsertOne(ctx, wallet)
	if err != nil {
		return err
	}
	wallet.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *walletRepo) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.wallets.FindOne(ctx, bson.M{"user_id": userID}).Decode(&wallet)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepo) UpdateBalance(ctx context.Context, userID primitive.ObjectID, balance uint) error {
	_, err := r.wallets.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"balance": balance}})
	return err
}

func (r *walletRepo) AddHistory(ctx context.Context, history *domain.WalletHistory) error {
	_, err := r.history.InsertOne(ctx, history)
	return err
}

func (r *walletRepo) FindHistoryByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.WalletHistory, error) {
	cursor, err := r.history.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []domain.WalletHistory
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, err
	}
	return histories, nil
}
