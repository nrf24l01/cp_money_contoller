package core

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Config struct {
	PGHost            string
	PGPort            string
	PGUser            string
	PGPassword        string
	PGDatabase        string
	PGSSLMode         string
	PGTimeZone        string

	RabbitMQHost      string
	RabbitMQPort      string
	RabbitMQUser      string
	RabbitMQPassword  string

	RedisHost         string
	RedisPort         string
	RedisPassword     string
	RedisDB           uint64
	KeysSetName       string
	 
	SyncInterval      time.Duration

	TaskQueue         string
}

func BuildConfigFromEnv() (*Config, error) {
	redis_db := os.Getenv("REDIS_DB")
	redisDB, err := strconv.ParseUint(redis_db, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REDIS_DB: %w", err)
	}

	syncIntervalStr := os.Getenv("SYNC_INTERVAL")
	syncInterval, err := parseDuration(syncIntervalStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SYNC_INTERVAL: %w", err)
	}

	cfg := &Config{
		PGHost:           os.Getenv("PG_HOST"),
		PGPort:           os.Getenv("PG_PORT"),
		PGUser:           os.Getenv("PG_USER"),
		PGPassword:       os.Getenv("PG_PASSWORD"),
		PGDatabase:       os.Getenv("PG_DATABASE"),
		PGSSLMode:        os.Getenv("PG_SSL_MODE"),
		PGTimeZone:       os.Getenv("PG_TIME_ZONE"),

		RabbitMQHost:     os.Getenv("RABBITMQ_HOST"),
		RabbitMQPort:     os.Getenv("RABBITMQ_PORT"),
		RabbitMQUser:     os.Getenv("RABBITMQ_USER"),
		RabbitMQPassword: os.Getenv("RABBITMQ_PASSWORD"),

		TaskQueue:        os.Getenv("TASK_QUEUE"),

		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:         redisDB,
		KeysSetName:     os.Getenv("REDIS_KEYS_SET_NAME"),
		SyncInterval:    syncInterval,
	}
	return cfg, nil
}
func (c *Config) GetAMQPURL() string {
	return "amqp://" + c.RabbitMQUser + ":" + c.RabbitMQPassword + "@" + c.RabbitMQHost + ":" + c.RabbitMQPort + "/"
}

func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	re := regexp.MustCompile(`(\d+)([dhms])`)
	matches := re.FindAllStringSubmatch(s, -1)
	var total time.Duration
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", match[1])
		}
		unit := match[2]
		switch unit {
		case "d":
			total += time.Duration(value) * 24 * time.Hour
		case "h":
			total += time.Duration(value) * time.Hour
		case "m":
			total += time.Duration(value) * time.Minute
		case "s":
			total += time.Duration(value) * time.Second
		default:
			return 0, fmt.Errorf("unknown unit: %s", unit)
		}
	}
	return total, nil
}