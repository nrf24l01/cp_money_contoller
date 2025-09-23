package rabbitmq

type Task struct {
	UUID      string       `json:"uuid"`
	Type      string       `json:"type"`
	Payload   interface{}  `json:"payload"`
	SecretKey string       `json:"secret_key"`
}