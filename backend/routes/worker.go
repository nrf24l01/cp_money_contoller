package routes

import (
	"github.com/nrf24l01/cp_money_controller/backend/handlers"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
	"github.com/nrf24l01/go-web-utils/echokit"

	"github.com/labstack/echo/v4"
)

func RegisterWorkerRoutes(e *echo.Echo, h *handlers.Handler) {
	group := e.Group("/worker")

	group.PUT("/task/:uuid", h.WorkerTaskUpdateHandler, echokit.ValidationMiddleware(func() interface{} {
		return &schemas.WorkerTaskUpdateRequest{}
	}))
}