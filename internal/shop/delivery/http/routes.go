package http

import (
	"github.com/22Fariz22/merch-shop/internal/middleware"
	"github.com/22Fariz22/merch-shop/internal/shop"
	"github.com/labstack/echo/v4"
)

// Map shop routes
func MapShopRoutes(shopGroup *echo.Group, h shop.Handlers, mw *middleware.MiddlewareManager) {
	shopGroup.GET("/:item", h.Buy(), mw.AuthJWTMiddleware())
}
