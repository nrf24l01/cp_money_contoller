package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_controller/backend/models"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
	"gorm.io/datatypes"
)

func (h *Handler) WorkerTaskUpdateHandler(c echo.Context) error {
	body := c.Get("validatedBody").(*schemas.WorkerTaskUpdateRequest)
	uuid := c.Param("uuid")

	var taskStatus models.TaskStatus
	if err := h.DB.Joins("Task").Where("Task.uuid = ?", uuid).First(&taskStatus).Error; err != nil {
		return c.JSON(http.StatusNotFound, schemas.DefaultNotFoundResponse)
	}

	taskStatus.Status = body.Status
	if jsonBytes, err := json.Marshal(body.Result); err != nil {
		return c.JSON(http.StatusBadRequest, schemas.DefaultBadRequestResponse)
	} else {
		taskStatus.Result = datatypes.JSON(jsonBytes)
	}
	if jsonBytes, err := json.Marshal(body.Logs); err != nil {
		return c.JSON(http.StatusBadRequest, schemas.DefaultBadRequestResponse)
	} else {
		taskStatus.Logs = datatypes.JSON(jsonBytes)
	}
	
	if err := h.DB.Save(&taskStatus).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	return c.JSON(http.StatusOK, schemas.DefaultSuccessResponse)
}
