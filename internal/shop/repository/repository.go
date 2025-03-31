package repository

import (
	"context"
	"fmt"

	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/shop"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

// shop Repository
type shopRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

// Shop repository constructor
func NewShopRepository(db *sqlx.DB, log logger.Logger) shop.Repository {
	return &shopRepo{db: db, logger: log}
}

func (r *shopRepo) Buy(ctx context.Context, user *models.User, item string, price int) error {
	fmt.Println("Here repo Buy()")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "authRepo.Buy.BeginTx")
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.logger.Errorf("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	// Получаем кошелёк
	wallet := &models.Wallet{}
	err = tx.QueryRowxContext(ctx, getWalletQuery, user.UserID).StructScan(wallet)
	if err != nil {
		return errors.Wrap(err, "authRepo.Buy.GetWallet")
	}

	if wallet.Balance < price {
		return errors.New("authRepo.Buy: insufficient balance")
	}

	// Обновляем баланс кошелька
	newBalance := wallet.Balance - price
	err = tx.QueryRowxContext(ctx, updateWalletQuery, newBalance, wallet.WalletID).StructScan(wallet)
	if err != nil {
		return errors.Wrap(err, "authRepo.Buy.UpdateWallet")
	}

	// Создаём запись о покупке
	purchase := &models.Purchase{}
	err = tx.QueryRowxContext(ctx, createPurchaseQuery, user.UserID, wallet.WalletID, item, price).
		StructScan(purchase)
	if err != nil {
		return errors.Wrap(err, "authRepo.Buy.CreatePurchase")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "authRepo.Buy.Commit")
	}

	r.logger.Debug("balance after purchase:", newBalance)
	return nil
}
