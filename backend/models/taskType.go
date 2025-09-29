package models

import (
	"github.com/nrf24l01/go-web-utils/goorm"
	"gorm.io/datatypes"
)

type TaskType struct {
	goorm.BaseModel
	Name	  string          `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Type      string          `json:"type" gorm:"type:varchar(255);not null"`
	Template  datatypes.JSON  `json:"template" gorm:"type:jsonb;not null"`
}