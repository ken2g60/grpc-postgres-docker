package models

import (
	"context"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UUID        string    `json:"uuid" gorm:"unique"`
	First_name  string    `json:"first_name" gorm:"index"`
	Last_name   string    `json:"last_name" gorm:"unique"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.UUID = uuid.New().String()
	return nil
}

func (user *User) BeforeSave(tx *gorm.DB) error {
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	return nil
}

func (user *User) AfterCreate(tx *gorm.DB) error {
	wallet := Wallet{
		UserID:           user.UUID,
		AvailableBalance: 0,
		Status:           "active",
		CreatedAt:        time.Now(),
	}
	tx.Model(&Wallet{}).Create(wallet)
	return nil
}

func CreateUser(ctx context.Context, db *gorm.DB, user *User) (err error) {
	err = db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func FindUser(ctx context.Context, db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.WithContext(ctx).Table("users").Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, nil
	}
	return &user, nil
}

func FindUserById(ctx context.Context, db *gorm.DB, user_id string) (*User, error) {
	var user User
	err := db.WithContext(ctx).Table("users").Where("uuid = ?", user_id).Find(&user).Error
	if err != nil {
		return nil, nil
	}
	return &user, nil
}
