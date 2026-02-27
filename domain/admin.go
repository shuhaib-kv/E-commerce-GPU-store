package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type PaymentGateway struct {
	Provider string `bson:"provider" json:"provider"` // e.g. "razorpay", "stripe"
	Key      string `bson:"key" json:"key"`
	Secret   string `bson:"secret" json:"-"`
}

type Admin struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Email          string             `bson:"email" json:"email"`
	Password       string             `bson:"password" json:"-"`
	Role           string             `bson:"role" json:"role"` // "admin" or "superadmin"
	StoreName      string             `bson:"store_name" json:"store_name"`
	BlockStatus    bool               `bson:"block_status" json:"block_status"`
	PaymentGateway PaymentGateway     `bson:"payment_gateway" json:"payment_gateway"`
}

func (a *Admin) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	a.Password = string(bytes)
	return nil
}

func (a *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

type AdminRepository interface {
	Create(ctx context.Context, admin *Admin) error
	FindByEmail(ctx context.Context, email string) (*Admin, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Admin, error)
	FindAll(ctx context.Context) ([]*Admin, error)
	Update(ctx context.Context, admin *Admin) error
}
