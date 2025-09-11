package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nrf24l01/cp_money_contoller/core"
	"github.com/nrf24l01/cp_money_contoller/models"
)

func (r *Redis) NewWorker(id uint64) (string, string, error) {
	worker := &models.Worker{
		UUID: uuid.New().String(),
		ID:   id,
		Key:  core.GenerateRandomString(32),
	}

	key := worker.UUID
	key = "worker:" + key
	value, err := json.Marshal(worker)
	if err != nil {
		return "", "", err
	}
	go func() {
		r.Client.Set(context.Background(), key, value, time.Duration(int64(r.WorkerLifetime))*time.Second).Err()
	}()
	return worker.UUID, worker.Key, nil
}
