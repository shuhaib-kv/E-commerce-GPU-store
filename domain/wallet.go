package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID  primitive.ObjectID `bson:"user_id" json:"user_id"`
	Balance uint               `bson:"balance" json:"balance"`
}

type WalletHistory struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Credit uint               `bson:"credit" json:"credit"`
	Debit  uint               `bson:"debit" json:"debit"`
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *Wallet) error
	FindByUserID(ctx context.Context, userID primitive.ObjectID) (*Wallet, error)
	UpdateBalance(ctx context.Context, userID primitive.ObjectID, balance uint) error
	AddHistory(ctx context.Context, history *WalletHistory) error
	FindHistoryByUserID(ctx context.Context, userID primitive.ObjectID) ([]WalletHistory, error)
}
