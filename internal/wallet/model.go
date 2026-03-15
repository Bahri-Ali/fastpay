package wallet

import (
	"time"
)

type Wallet struct {
    ID        string    `json:"id" db:"id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Balance   float64   `json:"balance" db:"balance"` // لتبسيط الأمور سنستخدم float64
    Currency  string    `json:"currency" db:"currency"`
    IsFrozen  bool      `json:"is_frozen" db:"is_frozen"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

