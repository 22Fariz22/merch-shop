package auth

import (
	"context"

	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/google/uuid"
)

// Auth usecase interface
type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
