package info

import (
	"context"

	"github.com/22Fariz22/merch-shop/internal/models"
)

// Info repository interface
type Repository interface {
	Info(ctx context.Context, user *models.User) (*models.Info, error)
}
