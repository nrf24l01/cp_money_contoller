package tasks

import (
	"errors"

	"github.com/nrf24l01/cp_money_controller/backend/models"
)

func (h *Handler) CreateUser(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	user := &models.User{
		Username: username,
	}
	user.SetPassword(password)

	if err := h.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}
