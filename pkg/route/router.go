package route

import "github.com/labstack/echo/v4"

type Router interface {
	Route(e *echo.Echo)
}