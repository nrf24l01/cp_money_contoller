package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_contoller/handlers"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler) {
	RegisterAuthRoutes(e, h)
}