package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/22Fariz22/merch-shop/internal/info"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// info Repository
type infoRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

// Info repository constructor
func NewInfoRepository(db *sqlx.DB, log logger.Logger) info.Repository {
	return &infoRepo{db: db, logger: log}
}

func (r *infoRepo) Info(ctx context.Context, user *models.User) (*models.Info, error) {
	infoRes := &models.Info{
		CoinHistory: models.CoinHistory{
			Received: []models.TransferInfo{},
			Sent:     []models.TransferInfo{},
		},
	}

	// Транзакция с уровнем изоляции Read Committed
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, errors.Wrap(err, "infoRepo.Info.BeginTx")
	}

	// Откат транзакции в случае ошибки
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
		if err == sql.ErrNoRows {
			return nil, errors.New("wallet not found")
		}
		return nil, errors.Wrap(err, "infoRepo.Info.GetWallet")
	}
	infoRes.Coins = wallet.Balance

	// Получаем историю переводов
	rows, err := tx.QueryxContext(ctx, getTransfers, user.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "infoRepo.Info.GetTransfers")
	}
	defer rows.Close()

	for rows.Next() {
		var t struct {
			TransferID   uuid.UUID `db:"transfer_id"`
			WalletFromID uuid.UUID `db:"wallet_from_id"`
			WalletToID   uuid.UUID `db:"wallet_to_id"`
			Amount       int       `db:"amount"`
			CreatedAt    time.Time `db:"created_at"`
			Status       string    `db:"status"`
			FromUsername *string   `db:"from_username"`
			ToUsername   *string   `db:"to_username"`
		}
		if err := rows.StructScan(&t); err != nil {
			return nil, errors.Wrap(err, "infoRepo.Info.ScanTransfer")
		}

		if t.WalletFromID == wallet.WalletID {
			if t.ToUsername != nil {
				infoRes.CoinHistory.Sent = append(infoRes.CoinHistory.Sent, models.TransferInfo{
					ToUser: *t.ToUsername,
					Amount: t.Amount,
				})
			}
		} else if t.WalletToID == wallet.WalletID {
			if t.FromUsername != nil {
				infoRes.CoinHistory.Received = append(infoRes.CoinHistory.Received, models.TransferInfo{
					FromUser: *t.FromUsername,
					Amount:   t.Amount,
				})
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "infoRepo.Info.GetTransfers.Rows")
	}

	// Получаем покупки
	purchaseRows, err := tx.QueryxContext(ctx, getPurchases, user.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "infoRepo.Info.GetPurchases")
	}
	defer purchaseRows.Close()

	for purchaseRows.Next() {
		var p models.PurchaseInfo
		if err := purchaseRows.StructScan(&p); err != nil {
			return nil, errors.Wrap(err, "infoRepo.Info.ScanPurchase")
		}
		infoRes.Inventory = append(infoRes.Inventory, p)
	}
	if err = purchaseRows.Err(); err != nil {
		return nil, errors.Wrap(err, "infoRepo.Info.GetPurchases.Rows")
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "infoRepo.Info.Commit")
	}

	return infoRes, nil
}
