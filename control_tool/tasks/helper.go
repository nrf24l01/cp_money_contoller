package tasks

import (
	"github.com/nrf24l01/cp_money_controller/control_tool/core"
	"gorm.io/gorm"
)

type Handler struct {
	DB      *gorm.DB
	Config  *core.Config
}