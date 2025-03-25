package utils

import (
	"context"
	"net/http"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/labstack/echo/v4"
)

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// UserCtxKey is a key used for the User object in the context
type UserCtxKey struct{}

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get user ip address
func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

// Get context  with request id
func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// Read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

// Error response with logging error for echo context
func LogResponseError(ctx echo.Context, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
}

// Configure jwt cookie
func CreateSessionCookie(cfg *config.Config, session string) *http.Cookie {
	return &http.Cookie{
		Name:  cfg.Session.Name,
		Value: session,
		Path:  "/",
		// Domain: "/",
		// Expires:    time.Now().Add(1 * time.Minute),
		RawExpires: "",
		MaxAge:     cfg.Session.Expire,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
	}
}

// Delete session
func DeleteSessionCookie(c echo.Context, sessionName string) {
	c.SetCookie(&http.Cookie{
		Name:   sessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
