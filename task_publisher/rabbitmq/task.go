package rabbitmq

type Task struct {
	Type      string
	Payload   interface{}
	SecretKey string
}