package main

import (
	"github.com/nrf24l01/cp_money_controller/backend/core"
	"github.com/nrf24l01/cp_money_controller/backend/handlers"
	"github.com/nrf24l01/cp_money_controller/backend/models"
	"github.com/nrf24l01/cp_money_controller/backend/routes"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
	"github.com/nrf24l01/go-web-utils/echokit"

	"github.com/go-playground/validator"

	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)
func main() {
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}
	
	config, err := core.BuildConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}

	// Data providers
	db := models.RegisterPostgres(config)

	// Validators
	validater := validator.New()

	e := echo.New()

	e.Validator = &echokit.CustomValidator{Validator: validater}

	if !config.ProductionEnv {
		e.Use(echoMw.Logger())
	}
    e.Use(echoMw.Recover())
	
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins: []string{config.AllowOrigins},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	log.Printf("CORS allowed origins: %s", config.AllowOrigins)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, schemas.Message{Status: "money controller is ok"})
	})

	handler := &handlers.Handler{DB: db, Config: config}
	routes.RegisterRoutes(e, handler)

	e.Logger.Fatal(e.Start(config.APPHost))
}