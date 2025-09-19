package models

import (
	"github.com/google/uuid"
	"github.com/nrf24l01/go-web-utils/goorm"
	"gorm.io/datatypes"
)

type Task struct {
	goorm.BaseModel
	Type      string         `gorm:"type:varchar(100);not null"`
	Payload   datatypes.JSON `gorm:"type:jsonb;not null"`
	SecretKey string         `gorm:"type:varchar(100)"`
}

type TaskStatus struct {
	goorm.BaseModel
	TaskID         uuid.UUID       `gorm:"column:task_id;not null"`
	Task           *Task           `gorm:"foreignKey:TaskID;references:ID"`
	Logs           datatypes.JSON  `gorm:"type:jsonb;default:'[]'"`
	Status         string          `gorm:"type:varchar(50);not null"`
	Result         datatypes.JSON  `gorm:"type:jsonb"`
}