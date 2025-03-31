package usecase

import (
	"context"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/info"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/pkg/logger"
)

// Info UseCase
type infoUC struct {
	cfg      *config.Config
	infoRepo info.Repository
	logger   logger.Logger
}

// Info UseCase constructor
func NewInfoUseCase(cfg *config.Config, infoRepo info.Repository, logger logger.Logger) info.UseCase {
	return &infoUC{cfg: cfg, infoRepo: infoRepo, logger: logger}
}

func (u *infoUC) Info(ctx context.Context, user *models.User) (*models.Info, error) {
	return u.infoRepo.Info(ctx, user)
}
