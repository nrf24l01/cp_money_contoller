package handlers

import (
	"gorm.io/gorm"

	"github.com/nrf24l01/cp_money_contoller/core"
)

type Handler struct {
	DB *gorm.DB
	Config *core.Config
}