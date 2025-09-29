package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_controller/backend/models"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
	"gorm.io/gorm"
)

func (h *Handler) AddTaskTypeHandler(c echo.Context) error {
	task_body := c.Get("validatedBody").(*schemas.CreateTaskTypeRequest)

	var exist models.TaskType
	if err := h.DB.Where("type = ?", task_body.TaskType).First(&exist).Error; err == nil {
		return c.JSON(http.StatusConflict, schemas.DefaultConflictResponse)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	b, err := json.Marshal(task_body.TaskTemplate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	task_type := &models.TaskType{
		Name:         task_body.TaskName,
		Type:         task_body.TaskType,
		Template:     b,
	}
	if err := h.DB.Create(&task_type).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	return c.JSON(http.StatusOK, schemas.CreateTaskTypeResponse{
		UUID:         task_type.ID.String(),
		TaskName:     task_type.Name,
		TaskType:     task_type.Type,
		TaskTemplate: task_body.TaskTemplate,
	})
}

func (h *Handler) ListTaskTypesHandler(c echo.Context) error {
	var task_types []models.TaskType
	if err := h.DB.Find(&task_types).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	var resp []schemas.TaskTypeWithUUID
	for _, t := range task_types {
		var template []schemas.TaskTypeField
		if err := json.Unmarshal(t.Template, &template); err != nil {
			return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
		}
		resp = append(resp, schemas.TaskTypeWithUUID{
			UUID: t.ID.String(),
			TaskType: schemas.TaskType{
				TaskName:     t.Name,
				TaskType:     t.Type,
				TaskTemplate: template,
			},
		})
	}

	return c.JSON(http.StatusOK, resp)
}