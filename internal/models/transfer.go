package models

import (
	"time"

	"github.com/google/uuid"
)

// Transfer model
type Transfer struct {
	TransferID   uuid.UUID `json:"transfer_id" gorm:"column:transfer_id;type:uuid;primaryKey;default:gen_random_uuid()" db:"transfer_id" redis:"transfer_id" validate:"required"`
	WalletFromID uuid.UUID `json:"wallet_from_id" gorm:"column:wallet_from_id;type:uuid;not null;index" db:"wallet_from_id" redis:"wallet_from_id" validate:"required"`
	WalletToID   uuid.UUID `json:"wallet_to_id" gorm:"column:wallet_to_id;type:uuid;not null;index" db:"wallet_to_id" redis:"wallet_to_id" validate:"required"`
	Amount       int       `json:"amount" gorm:"column:amount;not null;check:amount >= 0" db:"amount" redis:"amount" validate:"gte=0"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"column:created_at;default:current_timestamp" db:"created_at" redis:"created_at"`
	Status       string    `json:"status" gorm:"column:status;type:varchar(20);default:'completed'" db:"status" redis:"status" validate:"oneof=pending completed failed"`
	WalletFrom   Wallet    `gorm:"foreignKey:WalletFromID;references:WalletID;constraint:OnDelete:CASCADE"`
	WalletTo     Wallet    `gorm:"foreignKey:WalletToID;references:WalletID;constraint:OnDelete:CASCADE"`
}

type TransferRequest struct {
	ToUser string `json:"toUser" validate:"required,uuid"`
	Amount int    `json:"amount" validate:"required,gte=1"`
}
