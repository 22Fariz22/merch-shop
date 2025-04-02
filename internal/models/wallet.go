package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	WalletID  uuid.UUID `json:"wallet_id" db:"wallet_id" redis:"wallet_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" db:"user_id" redis:"user_id" validate:"required"`
	Balance   int       `json:"balance" db:"balance" redis:"balance" validate:"gte=0"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" redis:"updated_at"`
}
