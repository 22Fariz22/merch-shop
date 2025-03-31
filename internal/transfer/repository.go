package transfer

import (
	"context"

	"github.com/22Fariz22/merch-shop/internal/models"
)

type Repository interface {
	Transfer(ctx context.Context, user *models.User, toUser string, amount int) error
}
