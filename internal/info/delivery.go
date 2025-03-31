package info

import "github.com/labstack/echo/v4"

// Info HTTP Handlers interface
type Handlers interface {
	Info() echo.HandlerFunc
}
