package usecase

import (
	"context"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/transfer"
	"github.com/22Fariz22/merch-shop/pkg/logger"
)

// Transfer UseCase
type transferUC struct {
	cfg          *config.Config
	transferRepo transfer.Repository
	logger       logger.Logger
}

// Transfer UseCase constructor
func NewTransferUseCase(cfg *config.Config, transferRepo transfer.Repository, logger logger.Logger) transfer.UseCase {
	return &transferUC{cfg: cfg, transferRepo: transferRepo, logger: logger}
}

// Transfer item
func (u *transferUC) Transfer(ctx context.Context, user *models.User, toUser string, amount int) error {
	return u.transferRepo.Transfer(ctx, user, toUser, amount)
}
