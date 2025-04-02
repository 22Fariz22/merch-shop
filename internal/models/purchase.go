package models

import (
	"time"

	"github.com/google/uuid"
)

type Purchase struct {
	PurchaseID uuid.UUID `json:"purchase_id" db:"purchase_id" redis:"purchase_id" validate:"required"`
	UserID     uuid.UUID `json:"user_id" db:"user_id" redis:"user_id" validate:"required"`
	WalletID   uuid.UUID `json:"wallet_id" db:"wallet_id" redis:"wallet_id" validate:"required"`
	Item       string    `json:"item" db:"item" redis:"item" validate:"required"`
	Price      int       `json:"price" db:"price" redis:"price" validate:"gte=0"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
}
