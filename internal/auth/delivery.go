package auth

import "github.com/labstack/echo/v4"

// Auth HTTP Handlers interface
type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	// GetMe() echo.HandlerFunc
	// GetCSRFToken() echo.HandlerFunc
}
