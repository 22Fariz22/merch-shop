package models

import (
	"time"

	"github.com/google/uuid"
)

// Transfer model
type Transfer struct {
	TransferID   uuid.UUID `json:"transfer_id" db:"transfer_id" redis:"transfer_id" validate:"required"`
	WalletFromID uuid.UUID `json:"wallet_from_id" db:"wallet_from_id" redis:"wallet_from_id" validate:"required"`
	WalletToID   uuid.UUID `json:"wallet_to_id" db:"wallet_to_id" redis:"wallet_to_id" validate:"required"`
	Amount       int       `json:"amount" db:"amount" redis:"amount" validate:"gte=0"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
	Status       string    `json:"status" db:"status" redis:"status" validate:"oneof=pending completed failed"`
}

type TransferRequest struct {
	ToUser string `json:"toUser" validate:"required,uuid"`
	Amount int    `json:"amount" validate:"required,gte=1"`
}
