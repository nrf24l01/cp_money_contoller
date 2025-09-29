package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_controller/backend/handlers"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
	"github.com/nrf24l01/go-web-utils/echokit"
)

func RegisterTaskTypeRoutes(e *echo.Echo, h *handlers.Handler) {
	g := e.Group("/task/type")
	g.Use(echokit.JWTMiddleware([]byte(h.Config.JWTAccessSecret)))

	g.POST("", h.AddTaskTypeHandler, echokit.ValidationMiddleware(func() interface{} {
		return &schemas.CreateTaskTypeRequest{}
	}))

	g.GET("", h.ListTaskTypesHandler)
}