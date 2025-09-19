package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_controller/backend/models"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
	"gorm.io/datatypes"
)

func (h *Handler) CreateTaskHandler(c echo.Context) error {
	task_body := c.Get("validatedBody").(*schemas.CreateTaskRequest)

	// Ensure payload is stored as JSON bytes. The validator may leave Payload
	// as map[string]interface{} or []byte depending on middleware. Safely
	// marshal to []byte when necessary.
	var payloadBytes []byte
	switch v := task_body.Payload.(type) {
	case nil:
		payloadBytes = []byte("null")
	case []byte:
		payloadBytes = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
		}
		payloadBytes = b
	}

	task := &models.Task{
		Type:    task_body.Type,
		Payload: datatypes.JSON(payloadBytes),
	}

	if err := h.DB.Create(task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	return c.JSON(http.StatusOK, schemas.CreateTaskResponse{
		UUID: task.ID.String(),
	})
}

func (h *Handler) GetTaskHandler(c echo.Context) error {
	task_uuid := c.Param("uuid")

	var task models.Task
	if err := h.DB.First(&task, "id = ?", task_uuid).Error; err != nil {
		return c.JSON(http.StatusNotFound, schemas.DefaultNotFoundResponse)
	}

	var task_status models.TaskStatus
	var resp schemas.GetTaskResponse
	if err := h.DB.First(&task_status, "task_id = ?", task.ID).Error; err != nil {
		resp = schemas.GetTaskResponse{
			UUID:          task.ID.String(),
			Type:          task.Type,
			InputPayload:  task.Payload,
			OutputPayload: nil,
			Logs:          []string{},
			Status:        "pending",
			LastUpdate:    0,
		}
	} else {
		var logs []string
		if task_status.Logs != nil {
			json.Unmarshal(task_status.Logs, &logs)
		}

		resp = schemas.GetTaskResponse{
			UUID:          task.ID.String(),
			Type:          task.Type,
			InputPayload:  task.Payload,
			OutputPayload: task_status.Result,
			Logs:          logs,
			Status:        task_status.Status,
			LastUpdate:    uint64(task_status.UpdatedAt.Unix()),
		}
	}

	return c.JSON(http.StatusOK, resp)
}