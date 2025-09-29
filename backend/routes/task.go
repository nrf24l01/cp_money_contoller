package routes

import (
	"github.com/nrf24l01/cp_money_controller/backend/handlers"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
	"github.com/nrf24l01/go-web-utils/echokit"

	"github.com/labstack/echo/v4"
)

func RegisterTaskRoutes(e *echo.Echo, h *handlers.Handler) {
	group := e.Group("/task")
	// group.Use(echokit.JWTMiddleware([]byte(h.Config.JWTAccessSecret)))

	group.GET("/", h.GetTasksHandler)
	group.POST("/", h.CreateTaskHandler, echokit.ValidationMiddleware(func() interface{} {
		return &schemas.CreateTaskRequest{}
	}))
	group.GET("/:uuid", h.GetTaskHandler)
}