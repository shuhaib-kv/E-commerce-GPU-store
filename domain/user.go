package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName   string             `bson:"first_name" json:"first_name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	UserName    string             `bson:"user_name" json:"user_name"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"`
	Phone       string             `bson:"phone" json:"phone"`
	BlockStatus bool               `bson:"block_status" json:"block_status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Address struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Pincode     string             `bson:"pincode" json:"pincode"`
	House       string             `bson:"house" json:"house"`
	Area        string             `bson:"area" json:"area"`
	Landmark    string             `bson:"landmark" json:"landmark"`
	City        string             `bson:"city" json:"city"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindAll(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]User, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	CreateAddress(ctx context.Context, address *Address) error
	FindAddressByUserID(ctx context.Context, userID primitive.ObjectID) (*Address, error)
	UpdateAddress(ctx context.Context, userID primitive.ObjectID, update map[string]interface{}) error
}
