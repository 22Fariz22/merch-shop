package transfer

import "github.com/labstack/echo/v4"

type Handlers interface {
	Transfer() echo.HandlerFunc
}
