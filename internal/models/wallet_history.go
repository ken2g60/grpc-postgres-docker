package models

import "time"

type WalletHistory struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	WalletID  int       `json:"wallet_id"`
	UserID    string    `json:"user_id"`
	Available float32   `json:"available_balance"`
	CreatedAt time.Time `json:"created_at"`
}
