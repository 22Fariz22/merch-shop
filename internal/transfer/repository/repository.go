package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/transfer"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// transfer Repository
type transferRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

// Transfer repository constructor
func NewTransferRepository(db *sqlx.DB, log logger.Logger) transfer.Repository {
	return &transferRepo{db: db, logger: log}
}

const (
	// Запрос для получения кошелька с блокировкой строк
	getWalletQuery = `
        SELECT wallet_id, user_id, balance, updated_at 
        FROM wallets 
        WHERE user_id = $1 
        FOR UPDATE
    `
	// Запрос для обновления баланса кошелька
	updateWalletQuery = `
        UPDATE wallets 
        SET balance = $1, updated_at = NOW() 
        WHERE wallet_id = $2 
        RETURNING wallet_id, user_id, balance, updated_at
    `
	// Запрос для создания записи о переводе
	createTransferQuery = `
        INSERT INTO transfers (transfer_id, wallet_from_id, wallet_to_id, amount, created_at, status)
        VALUES ($1, $2, $3, $4, NOW(), 'completed')
        RETURNING transfer_id, wallet_from_id, wallet_to_id, amount, created_at, status
    `
	// Запрос для получения пользователя по UUID
	getUserByIDQuery = `
        SELECT user_id, username, password, created_at 
        FROM users 
        WHERE user_id = $1
    `
)

func (r *transferRepo) Transfer(ctx context.Context, user *models.User, toUser string, amount int) error {
	fmt.Println("Here repo Transfer()")
	r.logger.Debugf("Starting transfer from %s to %s for amount %d", user.Username, toUser, amount)

	// Проверка входных данных
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Парсим toUser как UUID
	toUserID, err := uuid.Parse(toUser)
	if err != nil {
		return errors.Wrap(err, "invalid toUser UUID")
	}

	// Начинаем транзакцию с уровнем изоляции SERIALIZABLE
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return errors.Wrap(err, "transferRepo.Transfer.BeginTx")
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.logger.Errorf("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	// Получаем кошелёк отправителя с блокировкой
	walletFrom := &models.Wallet{}
	err = tx.QueryRowxContext(ctx, getWalletQuery, user.UserID).StructScan(walletFrom)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("sender wallet not found")
		}
		return errors.Wrap(err, "transferRepo.Transfer.GetWalletFrom")
	}

	// Проверяем баланс отправителя
	if walletFrom.Balance < amount {
		return errors.New("insufficient balance")
	}

	// Получаем пользователя получателя по UUID
	recipient := &models.User{}
	err = tx.QueryRowxContext(ctx, getUserByIDQuery, toUserID).StructScan(recipient)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("recipient user not found")
		}
		return errors.Wrap(err, "transferRepo.Transfer.GetRecipient")
	}

	// Получаем кошелёк получателя с блокировкой
	walletTo := &models.Wallet{}
	err = tx.QueryRowxContext(ctx, getWalletQuery, recipient.UserID).StructScan(walletTo)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("recipient wallet not found")
		}
		return errors.Wrap(err, "transferRepo.Transfer.GetWalletTo")
	}

	// Обновляем баланс отправителя
	newBalanceFrom := walletFrom.Balance - amount
	err = tx.QueryRowxContext(ctx, updateWalletQuery, newBalanceFrom, walletFrom.WalletID).StructScan(walletFrom)
	if err != nil {
		return errors.Wrap(err, "transferRepo.Transfer.UpdateWalletFrom")
	}

	// Обновляем баланс получателя
	newBalanceTo := walletTo.Balance + amount
	err = tx.QueryRowxContext(ctx, updateWalletQuery, newBalanceTo, walletTo.WalletID).StructScan(walletTo)
	if err != nil {
		return errors.Wrap(err, "transferRepo.Transfer.UpdateWalletTo")
	}

	// Создаём запись о переводе
	transferID := uuid.New() // Генерируем уникальный transfer_id
	transfer := &models.Transfer{}
	err = tx.QueryRowxContext(ctx, createTransferQuery, transferID, walletFrom.WalletID, walletTo.WalletID, amount).
		StructScan(transfer)
	if err != nil {
		return errors.Wrap(err, "transferRepo.Transfer.CreateTransfer")
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "transferRepo.Transfer.Commit")
	}

	r.logger.Infof("Transfer completed: ID=%s, From=%s, To=%s, Amount=%d", transfer.TransferID, user.Username, toUser, amount)
	return nil
}
