package schemas

type WorkerTaskUpdateRequest struct {
	Result interface{} `json:"result" validate:"required"`
	Status string      `json:"status" validate:"required"`
	Logs   []string    `json:"logs" validate:"required"`
}