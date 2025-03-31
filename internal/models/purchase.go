package models

import (
	"time"

	"github.com/google/uuid"
)

type Purchase struct {
	PurchaseID uuid.UUID `json:"purchase_id" gorm:"column:purchase_id;type:uuid;primaryKey;default:gen_random_uuid()" db:"purchase_id" redis:"purchase_id" validate:"required"`
	UserID     uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null;index" db:"user_id" redis:"user_id" validate:"required"`
	WalletID   uuid.UUID `json:"wallet_id" gorm:"column:wallet_id;type:uuid;not null;index" db:"wallet_id" redis:"wallet_id" validate:"required"`
	Item       string    `json:"item" gorm:"column:item;check:item != ''" db:"item" redis:"item" validate:"required"`
	Price      int       `json:"price" gorm:"column:price;not null;check:price >= 0" db:"price" redis:"price" validate:"gte=0"`
	CreatedAt  time.Time `json:"created_at,omitempty" gorm:"column:created_at;default:current_timestamp" db:"created_at" redis:"created_at"`
	User       User      `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
	Wallet     Wallet    `gorm:"foreignKey:WalletID;references:WalletID;constraint:OnDelete:CASCADE"`
}
