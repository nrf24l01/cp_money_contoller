package schemas

type CreateTaskRequest struct {
	Type 	  string      `json:"type" validate:"required"`
	Payload   interface{} `json:"payload"`
}

type CreateTaskResponse struct {
	UUID string `json:"uuid"`
}

type GetTaskResponse struct {
	UUID           string      `json:"uuid"`
	Type           string      `json:"type"`
	InputPayload   interface{} `json:"payload_input"`
	OutputPayload  interface{} `json:"payload_output"`
	Logs           []string    `json:"logs"`
	Status         string      `json:"status"`
	LastUpdate     uint64      `json:"last_update"`
}