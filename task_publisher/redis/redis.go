package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nrf24l01/cp_money_controller/task_publisher/core"
	"github.com/redis/go-redis/v9"
)

// RedisClient обёртка для Redis
type RedisClient struct {
    client       *redis.Client
    ctx          context.Context
    keysSetName  string
}

func CreateRedisFromCFG(cfg *core.Config) (*RedisClient, error) {
	return NewRedisClient(
		fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		cfg.RedisPassword,
		int(cfg.RedisDB),
		cfg.KeysSetName,
	), nil
}

// NewRedisClient создаёт подключение к Redis
func NewRedisClient(addr string, password string, db int, keysSetName string) *RedisClient {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    ctx := context.Background()

    // Проверка соединения
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Не удалось подключиться к Redis: %v", err)
    }

    return &RedisClient{
        client:      rdb,
        ctx:         ctx,
        keysSetName: keysSetName,
    }
}

// Set добавляет ключ-значение с опциональным временем жизни
func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
    return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get возвращает значение по ключу
func (r *RedisClient) Get(key string) (string, error) {
    val, err := r.client.Get(r.ctx, key).Result()
    if err == redis.Nil {
        return "", fmt.Errorf("ключ %s не найден", key)
    }
    return val, err
}

func (r *RedisClient) Close() error {
    return r.client.Close()
}

func (r *RedisClient) AddIdToSet(id string) error {
    return r.client.SAdd(r.ctx, r.keysSetName, id).Err()
}

func (r *RedisClient) GetIDSFromSet() ([]string, error) {
    return r.client.SMembers(r.ctx, r.keysSetName).Result()
}

func (r *RedisClient) Purge() error {
    return r.client.FlushDB(r.ctx).Err()
}