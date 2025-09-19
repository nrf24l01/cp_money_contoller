package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nrf24l01/cp_money_controller/control_tool/core"
	"github.com/nrf24l01/cp_money_controller/control_tool/database"
	"github.com/nrf24l01/cp_money_controller/control_tool/tasks"
)

func main() {
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	cfg := core.BuildConfigFromEnv()
	db := database.RegisterPostgres(cfg)

	h := tasks.Handler{
		DB: db,
		Config: cfg,
	}

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 3 && argsWithoutProg[0] == "create-user" {
		user, password := argsWithoutProg[1], argsWithoutProg[2]
		if err := h.CreateUser(user, password); err != nil {
			panic(err)
		}
		log.Printf("User %s created successfully", user)
	}
}