package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_controller/backend/core"
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
		SecretKey: core.GenerateRandomString(100),
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

func (h *Handler) GetTasksHandler(c echo.Context) error {
	var tasks []models.Task
	if err := h.DB.Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	var statuses []models.TaskStatus
	if err := h.DB.Find(&statuses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	statusMap := make(map[string]models.TaskStatus)
	for _, ts := range statuses {
		statusMap[ts.TaskID.String()] = ts
	}

	var responses []schemas.GetTaskResponse
	for _, task := range tasks {
		var logs []string
		var status string
		var result interface{}
		var lastUpdate uint64

		if ts, ok := statusMap[task.ID.String()]; ok {
			if ts.Logs != nil {
				json.Unmarshal(ts.Logs, &logs)
			}
			status = ts.Status
			result = ts.Result
			lastUpdate = uint64(ts.UpdatedAt.Unix())
		} else {
			status = "pending"
			result = nil
			logs = []string{}
			lastUpdate = 0
		}

		resp := schemas.GetTaskResponse{
			UUID:          task.ID.String(),
			Type:          task.Type,
			InputPayload:  task.Payload,
			OutputPayload: result,
			Logs:          logs,
			Status:        status,
			LastUpdate:    lastUpdate,
		}
		responses = append(responses, resp)
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *Handler) GetLogsHandler(c echo.Context) error {
	task_uuid := c.Param("uuid")

	var task models.Task
	if err := h.DB.First(&task, "id = ?", task_uuid).Error; err != nil {
		return c.JSON(http.StatusNotFound, schemas.DefaultNotFoundResponse)
	}

	var task_status models.TaskStatus
	var logs []string
	if err := h.DB.First(&task_status, "task_id = ?", task.ID).Error; err != nil {
		// No status found, return empty logs
		logs = []string{}
	} else {
		if task_status.Logs != nil {
			json.Unmarshal(task_status.Logs, &logs)
		}
	}

	return c.JSON(http.StatusOK, logs)
}