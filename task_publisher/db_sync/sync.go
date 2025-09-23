package dbsync

import (
	"log"

	"github.com/nrf24l01/cp_money_controller/backend/models"
	"github.com/nrf24l01/cp_money_controller/task_publisher/rabbitmq"
)


func (h *Handler) SyncTasks() {
    log.Printf("Starting task synchronization...")

	// Get messages which are already in RabbitMQ from Redis
    log.Printf("Getting tasks which are in RabbitMQ from redis...")
    taskKeys, err := h.Redis.GetIDSFromSet()
    if err != nil {
        log.Printf("Error getting task keys from Redis: %v", err)
        return
    }
    log.Printf("Found %d tasks in Redis", len(taskKeys))

	// Get tasks from DB which are not in RabbitMQ
    var tasks []models.Task
    db := h.DB
    if len(taskKeys) > 0 { // If there are keys in Redis, filter them out
        db = db.Where("tasks.id NOT IN ?", taskKeys)
    }
    err = db.
        Joins("LEFT JOIN task_statuses ON tasks.id = task_statuses.task_id").
        Where("task_statuses.task_id IS NULL").
        Find(&tasks).Error
    if err != nil {
        log.Printf("Error querying tasks: %v", err)
        return
    }
    log.Printf("Found %d tasks to sync", len(tasks))

	// Check if no tasks to sync
	if len(tasks) == 0 {
		log.Printf("No tasks to sync. Exiting.")
		return
	}

	// Sync tasks to RabbitMQ and Redis
	log.Printf("Syncing tasks to RabbitMQ and Redis...")
	for _, task := range tasks {
		rmq_task := rabbitmq.Task{
			UUID:      task.ID.String(),
			Type:      task.Type,
			Payload:   task.Payload,
			SecretKey: task.SecretKey,
		}
		err := h.RMQ.AddTask(rmq_task)
		if err != nil {
			log.Printf("Error adding task to RabbitMQ: %v", err)
			continue
		}
		err = h.Redis.AddIdToSet(task.ID.String())
		if err != nil {
			log.Printf("Error adding task ID to Redis: %v", err)
			continue
		}
	}
	log.Printf("Task synchronization completed.")
}