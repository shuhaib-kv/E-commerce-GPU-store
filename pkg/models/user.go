package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID             uint      `json:"id" gorm:"primaryKey;unique"  `
	FirstName      string    `json:"fname"`
	LastName       string    `json:"lname"`
	UserName       string    `json:"uname"`
	Email          string    `gorm:"unique"`
	Password       string    `json:"password"`
	Phone          string    `json:"phone"`
	AddressId      uint      `json:"addressid"`
	Tocken         string    `json:"tocken"`
	Refresh_tocken string    `json:"refresh_tocken"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
	User_id        string    `json:"user_id"`
	Block_status   bool      `json:"block_status"`
}

func (user *Users) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *Users) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

type Address struct {
	gorm.Model

	UserId       uint   `json:"user_id"  gorm:"not null" `
	Name         string `json:"name"  gorm:"not null" `
	Phone_number uint   `json:"phone_number"  gorm:"not null" `
	Pincode      uint   `json:"pincode"  gorm:"not null" `
	House        string `json:"house"   `
	Area         string `json:"area"   `
	Landmark     string `json:"landmark"  gorm:"not null" `
	City         string `json:"city"  gorm:"not null" `
}
