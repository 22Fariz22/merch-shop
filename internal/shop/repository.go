package shop

import (
	"context"

	"github.com/22Fariz22/merch-shop/internal/models"
)

type Repository interface {
	Buy(ctx context.Context, user *models.User, item string, price int) error
}
