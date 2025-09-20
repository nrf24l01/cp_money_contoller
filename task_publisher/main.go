package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nrf24l01/cp_money_controller/task_publisher/core"
	"github.com/nrf24l01/cp_money_controller/task_publisher/database"
	dbsync "github.com/nrf24l01/cp_money_controller/task_publisher/db_sync"
	"github.com/nrf24l01/cp_money_controller/task_publisher/rabbitmq"
	"github.com/nrf24l01/cp_money_controller/task_publisher/redis"
)

func main() {
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	cfg, err := core.BuildConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to build config from env: %v", err)
		os.Exit(1)
	}

	db := database.RegisterPostgres(cfg)
	queue, err := rabbitmq.RabbitMQQueueFromCFG(cfg)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ queue: %v", err)
		os.Exit(1)
	}
	defer queue.Close()
	redis, err := redis.CreateRedisFromCFG(cfg)
	if err != nil {
		log.Fatalf("failed to create Redis client: %v", err)
		os.Exit(1)
	}

	h := dbsync.Handler{
		DB:    db,
		Redis: redis,
		RMQ:  queue,
		Cfg:   cfg,
	}

	h.StartSyncing()
}