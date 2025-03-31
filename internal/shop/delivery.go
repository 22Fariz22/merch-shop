package shop

import "github.com/labstack/echo/v4"

type Handlers interface {
	Buy() echo.HandlerFunc
}
