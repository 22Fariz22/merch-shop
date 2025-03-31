package http

import (
	"github.com/22Fariz22/merch-shop/internal/info"
	"github.com/22Fariz22/merch-shop/internal/middleware"
	"github.com/labstack/echo/v4"
)

// Map info routes
func MapInfoRoutes(infoGroup *echo.Group, h info.Handlers, mw *middleware.MiddlewareManager) {
	infoGroup.GET("", h.Info(), mw.AuthJWTMiddleware())
}
