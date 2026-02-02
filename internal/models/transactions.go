package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	UUID          string    `json:"uuid" gorm:"unique"`
	UserID        string    `json:"user_id"`
	Amount        float32   `json:"amount"`
	Description   string    `json:"description"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (transaction *Transaction) BeforeCreate(tx *gorm.DB) error {
	transaction.UUID = uuid.New().String()
	return nil
}

func CreateTransaction(ctx context.Context, db *gorm.DB, Transaction *Transaction) (err error) {
	err = db.WithContext(ctx).Create(&Transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func TransactionHistory(ctx context.Context, db *gorm.DB, user_id string) (*[]Transaction, error) {
	var transaction []Transaction
	err := db.WithContext(ctx).Table("transactions").Where("uuid = ?", user_id).Find(&transaction).Error
	if err != nil {
		return nil, nil
	}
	return &transaction, nil
}
