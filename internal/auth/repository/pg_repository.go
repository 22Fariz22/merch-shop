package repository

import (
	"context"

	"github.com/22Fariz22/merch-shop/internal/auth"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Auth Repository
type authRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

// Auth Repository constructor
func NewAuthRepository(db *sqlx.DB, log logger.Logger) auth.Repository {
	return &authRepo{db: db, logger: log}
}

// Create new user and wallet
func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	r.logger.Debug("here repo Register")

	// Начинаем транзакцию
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.BeginTx")
	}

	// Откладываем rollback на случай ошибки
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.logger.Errorf("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	createdUser := &models.User{}

	if err := tx.QueryRowxContext(ctx, createUserQuery, &user.Username, &user.Password).StructScan(createdUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.StructScan")
	}

	// Создаём кошелёк для пользователя
	initialBalance := 1000
	createdWallet := &models.Wallet{}
	err = tx.QueryRowxContext(ctx, createWalletQuery, createdUser.UserID, initialBalance).
		StructScan(createdWallet)
	if err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.CreateWallet.StructScan")
	}

	// Если всё успешно, коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.Commit")
	}

	r.logger.Debugf("User registered with ID: %s, Wallet ID: %s", createdUser.UserID, createdWallet.WalletID)
	return createdUser, nil
}

// Get user by id
func (r *authRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	r.logger.Debug("here repo GetByID")
	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserQuery, userID).StructScan(user); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetByID.QueryRowxContext")
	}
	return user, nil
}

// Find user by username
func (r *authRepo) FindByUsername(ctx context.Context, user *models.User) (*models.User, error) {
	r.logger.Debug("here repo FindByUsername")

	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByUsername, user.Username).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByUsername.QueryRowxContext")
	}
	return foundUser, nil
}
