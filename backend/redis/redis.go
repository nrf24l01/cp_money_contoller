package redis

import (
	"github.com/nrf24l01/cp_money_contoller/core"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
	WorkerLifetime uint64
}

func NewRedisClient(cfg *core.Config) Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return Redis{Client: rdb, WorkerLifetime: cfg.WorkerLifetime}
}

func (r *Redis) Close() error {
	return r.Client.Close()
}