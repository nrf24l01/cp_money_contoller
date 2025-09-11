package core

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	PGHost        string
	PGPort        string
	PGUser        string
	PGPassword    string
	PGDatabase    string
	PGSSLMode     string
	PGTimeZone    string

	WorkerLifetime uint64 // in seconds
}

func BuildConfigFromEnv() (*Config, error) {
	workerLifetime, err := strconv.ParseUint(os.Getenv("WORKER_LIFETIME"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid WORKER_LIFETIME: %v", err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %v", err)
	}

	cfg := &Config{
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RedisDB:        redisDB,

		PGHost:         os.Getenv("PG_HOST"),
		PGPort:         os.Getenv("PG_PORT"),
		PGUser:         os.Getenv("PG_USER"),
		PGPassword:     os.Getenv("PG_PASSWORD"),
		PGDatabase:     os.Getenv("PG_DATABASE"),
		PGSSLMode:      os.Getenv("PG_SSL_MODE"),
		PGTimeZone:     os.Getenv("PG_TIME_ZONE"),
		WorkerLifetime: workerLifetime,
	}
	return cfg, nil
}
