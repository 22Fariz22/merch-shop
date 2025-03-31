package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/auth"
	"github.com/22Fariz22/merch-shop/internal/models"

	"github.com/22Fariz22/merch-shop/pkg/httpErrors"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/22Fariz22/merch-shop/pkg/utils"
	"github.com/google/uuid"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	logger    logger.Logger
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo auth.RedisRepository, log logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, logger: log}
}

// Create new user
func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {

	existsUser, err := u.authRepo.FindByUsername(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrUsernameAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	// Кешируем пользователя
	err = u.redisRepo.SetUserCtx(ctx, createdUser.UserID.String(), 3600, createdUser)
	if err != nil {
		u.logger.Errorf("Failed to cache user: %v", err)
	}

	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Get user by id
func (u *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID.String()))
	if err != nil {
		u.logger.Errorf("authUC.GetByID.GetByIDCtx: %v", err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(userID.String()), cacheDuration, user); err != nil {
		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
	}

	user.SanitizePassword()

	return user, nil
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	u.logger.Debug("Here UC Login()")

	foundUser, err := u.authRepo.FindByUsername(ctx, user)
	if err != nil {
		u.logger.Debug("Here UC Login() FindByUsername err:", err)
		//если не нашли такого юзера, то регистрируем
		//TODO:нужно сравнить ошибку чтобы понимать ошибка БД или юзер не найден
		if err = user.PrepareCreate(); err != nil {
			return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
		}

		foundUser, err = u.authRepo.Register(ctx, user)
		if err != nil {
			u.logger.Debug("Here UC Login() Register err:", err)
			return nil, err
		}
	} else {
		//если нашли такого юзера, то проверяем пароль
		if err = foundUser.ComparePasswords(user.Password); err != nil {
			u.logger.Debug("Here UC Login() ComparePasswords err:", err)
			return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePasswords"))
		}
		//TODO: удаляем его текущую сессию в редисе?
	}

	foundUser.SanitizePassword()

	//генерируем токен
	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		u.logger.Debug("Here UC Login() GenerateJWTToken err:", err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.GetUsers.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *authUC) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
