package usecase

import (
	"context"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/shop"
	"github.com/22Fariz22/merch-shop/pkg/logger"
)

// Shop UseCase
type shopUC struct {
	cfg      *config.Config
	shopRepo shop.Repository
	logger   logger.Logger
}

// Shop UseCase constructor
func NewShopUseCase(cfg *config.Config, shopRepo shop.Repository, logger logger.Logger) shop.UseCase {
	return &shopUC{cfg: cfg, shopRepo: shopRepo, logger: logger}
}

// Buy item
func (u *shopUC) Buy(ctx context.Context, user *models.User, item string, price int) error {
	return u.shopRepo.Buy(ctx, user, item, price)
}
