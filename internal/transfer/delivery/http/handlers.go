package http

import (
	"net/http"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/transfer"
	"github.com/22Fariz22/merch-shop/pkg/httpErrors"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/22Fariz22/merch-shop/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Transfer handlers
type transferHandlers struct {
	cfg        *config.Config
	transferUC transfer.UseCase
	logger     logger.Logger
}

// NewTransferHandlers Transfer handlers constructor
func NewTransferHandlers(cfg *config.Config, transferUC transfer.UseCase, logger logger.Logger) transfer.Handlers {
	return &transferHandlers{cfg: cfg, transferUC: transferUC, logger: logger}
}

func (h *transferHandlers) Transfer() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debug("Here in Transfer handler")
		ctx := utils.GetRequestCtx(c)

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		h.logger.Debug("user:", user)

		var transferRequest models.TransferRequest
		if err := c.Bind(&transferRequest); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		}

		// Валидация запроса
		validate := validator.New()
		if err := validate.Struct(transferRequest); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed: " + err.Error()})
		}

		h.logger.Debug("transferRequests:", transferRequest)

		//usecase transfer
		err = h.transferUC.Transfer(ctx, user, transferRequest.ToUser, transferRequest.Amount)
		if err != nil {

			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		h.logger.Infof("Successfully sent %d coins from %s to %s", transferRequest.Amount, user.Username, transferRequest.ToUser)
		return c.JSON(http.StatusOK, map[string]string{"message": "Coins sent successfully"})
	}
}
