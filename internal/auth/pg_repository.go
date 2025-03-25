package auth

import (
	"context"

	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/google/uuid"
)

// Auth repository interface
type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByUsername(ctx context.Context, user *models.User) (*models.User, error)
		Delete(ctx context.Context, userID uuid.UUID) error
}
