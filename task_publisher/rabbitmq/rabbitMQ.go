package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/nrf24l01/cp_money_controller/task_publisher/core"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitQueue struct {
    Conn    *amqp.Connection
    Channel *amqp.Channel
    Queue   string
}

func RabbitMQQueueFromCFG(cfg *core.Config) (*RabbitQueue, error) {
	return NewRabbitQueue(cfg.GetAMQPURL(), cfg.TaskQueue)
}

func NewRabbitQueue(amqpURL, queueName string) (*RabbitQueue, error) {
    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }

    _, err = ch.QueueDeclare(
        queueName,
        false,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, fmt.Errorf("failed to declare queue: %w", err)
    }

    return &RabbitQueue{
        Conn:    conn,
        Channel: ch,
        Queue:   queueName,
    }, nil
}

func (r *RabbitQueue) Purge() error {
    _, err := r.Channel.QueuePurge(r.Queue, false)
    if err != nil {
        return fmt.Errorf("failed to purge queue: %w", err)
    }
    return nil
}

func (r *RabbitQueue) AddTask(task Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}
	err = r.Channel.Publish(
		"",       // exchange
		r.Queue,  // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
    if err != nil {
        return fmt.Errorf("failed to publish task: %w", err)
    }
    return nil
}

func (r *RabbitQueue) Close() {
    r.Channel.Close()
    r.Conn.Close()
}
