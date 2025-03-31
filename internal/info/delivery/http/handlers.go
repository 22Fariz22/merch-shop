package http

import (
	"net/http"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/info"
	"github.com/22Fariz22/merch-shop/pkg/httpErrors"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/22Fariz22/merch-shop/pkg/utils"
	"github.com/labstack/echo/v4"
)

// Info handlers
type infoHandlers struct {
	cfg    *config.Config
	infoUC info.UseCase
	logger logger.Logger
}

// NewInfoHandlers Info handlers constructor
func NewHandlers(cfg *config.Config, infoUC info.UseCase, logger logger.Logger) info.Handlers {
	return &infoHandlers{cfg: cfg, infoUC: infoUC, logger: logger}
}

func (h *infoHandlers) Info() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debug("Here in Info handler")
		ctx := utils.GetRequestCtx(c)

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		h.logger.Debug("user:", user)

		info, err := h.infoUC.Info(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		h.logger.Debug("Info:", info)

		return c.JSON(http.StatusOK, info)
	}
}
