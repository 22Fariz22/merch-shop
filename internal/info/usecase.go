package info

import (
	"context"
	"github.com/22Fariz22/merch-shop/internal/models"
)

type UseCase interface {
	Info(ctx context.Context, user *models.User) (*models.Info, error)
}
