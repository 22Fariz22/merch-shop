package http

import (
	"net/http"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/auth"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/session"
	"github.com/22Fariz22/merch-shop/pkg/httpErrors"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/22Fariz22/merch-shop/pkg/utils"
	"github.com/labstack/echo/v4"
)

// Auth handlers
type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	sessUC session.UCSession
	logger logger.Logger
}

// NewAuthHandlers Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, sessUC session.UCSession, log logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, sessUC: sessUC, logger: log}
}

// Если есть логин и пароль идем проверять в бд, если такой есть --выдаем токен,
// если нету --регимся и выдаем токен
func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Username string `json:"username" db:"username" validate:"required,lte=60"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	return func(c echo.Context) error {
		h.logger.Debug("Here handler Login()")

		ctx := utils.GetRequestCtx(c)

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			h.logger.Debug("Here handler Login() ReadRequest:", err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		//идем проверять юзера и получаем токен
		userWithToken, err := h.authUC.Login(ctx, &models.User{
			Username: login.Username,
			Password: login.Password,
		})
		if err != nil {
			h.logger.Debug("Here handler Login() h.authUC.Login err:", err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		//кэшируем userID и срок жизни в редисе
		sess, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: userWithToken.User.UserID,
		}, h.cfg.Session.Expire)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, sess))

		return c.JSON(http.StatusOK, userWithToken)
	}
}
