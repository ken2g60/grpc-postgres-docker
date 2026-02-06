package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	UUID             string    `json:"uuid" gorm:"unique"`
	UserID           string    `json:"user_id"`
	Country          string    `json:"country"`
	AvailableBalance float32   `json:"available_balance"`
	Status           string    `json:"status" gorm:"default:'active'"` // active, locked, unlocked, suspended
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (balance *Wallet) BeforeCreate(tx *gorm.DB) error {
	balance.UUID = uuid.New().String()
	return nil
}

func CreateWalletBalance(ctx context.Context, db *gorm.DB, Wallet *Wallet) (err error) {
	err = db.WithContext(ctx).Create(&Wallet).Error
	if err != nil {
		return err
	}
	return nil
}

// GetWalletByUserID retrieves a wallet by user ID
func GetWalletByUserID(ctx context.Context, db *gorm.DB, userID string) (*Wallet, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	var wallet Wallet
	result := db.Where("user_id = ?", userID).Find(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}

	return &wallet, nil
}

func (wallet *Wallet) Save(db *gorm.DB) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	// Start transaction
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.UpdateColumns(wallet)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	history := WalletHistory{
		WalletID:  wallet.ID,
		UserID:    wallet.UserID,
		Available: wallet.AvailableBalance,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
