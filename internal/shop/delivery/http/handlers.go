package http

import (
	"net/http"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/shop"
	"github.com/22Fariz22/merch-shop/pkg/httpErrors"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/22Fariz22/merch-shop/pkg/utils"
	"github.com/labstack/echo/v4"
)

// Shop handlers
type shopHandlers struct {
	cfg    *config.Config
	shopUC shop.UseCase
	logger logger.Logger
}

// NewShopHandlers Shop handlers constructor
func NewShopHandlers(cfg *config.Config, shopUC shop.UseCase, logger logger.Logger) shop.Handlers {
	return &shopHandlers{cfg: cfg, shopUC: shopUC, logger: logger}
}

func (h *shopHandlers) Buy() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debug("Here in buy handler")
		ctx := utils.GetRequestCtx(c)

		user, err := utils.GetUserFromCtx(ctx)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		h.logger.Debug("user:", user)

		item := c.Param("item")
		price, exists := models.DefaultShop.GetPrice(item)
		if !exists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Item not found"})
		}

		//byu and return balance,error
		err = h.shopUC.Buy(ctx, user, item, price)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Item bought successfully"})
	}
}
