package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_controller/backend/handlers"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler) {
	RegisterAuthRoutes(e, h)
	RegisterTaskRoutes(e, h)
	RegisterWorkerRoutes(e, h)
	RegisterAuthRoutes(e, h)
}