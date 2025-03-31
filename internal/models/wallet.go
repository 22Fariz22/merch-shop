package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	WalletID  uuid.UUID `json:"wallet_id" gorm:"column:wallet_id;type:uuid;primaryKey;default:gen_random_uuid()" db:"wallet_id" redis:"wallet_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid;not null;index" db:"user_id" redis:"user_id" validate:"required"`
	Balance   int       `json:"balance" gorm:"column:balance;default:1000;check:balance >= 0" db:"balance" redis:"balance" validate:"gte=0"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;default:current_timestamp" db:"updated_at" redis:"updated_at"`
	User      User      `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
}
