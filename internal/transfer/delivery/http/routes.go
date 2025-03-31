package http

import (
	"github.com/22Fariz22/merch-shop/internal/middleware"
	"github.com/22Fariz22/merch-shop/internal/transfer"
	"github.com/labstack/echo/v4"
)

// Map transfer routes
func MapTransferRoutes(transferGroup *echo.Group, h transfer.Handlers, mw *middleware.MiddlewareManager) {
	transferGroup.POST("", h.Transfer(), mw.AuthJWTMiddleware())
}
