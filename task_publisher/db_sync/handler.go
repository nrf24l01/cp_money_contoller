package dbsync

import (
	"github.com/nrf24l01/cp_money_controller/task_publisher/core"
	"github.com/nrf24l01/cp_money_controller/task_publisher/rabbitmq"
	"github.com/nrf24l01/cp_money_controller/task_publisher/redis"
	"gorm.io/gorm"
)

type Handler struct {
	Cfg   *core.Config
	Redis *redis.RedisClient
	RMQ   *rabbitmq.RabbitQueue
	DB    *gorm.DB
}