package models

import (
	"github.com/nrf24l01/go-web-utils/goorm"
	"github.com/nrf24l01/go-web-utils/passhash"
)

type User struct {
	goorm.BaseModel
	Username     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
}

func (u *User) SetPassword(password string) error {
	var err error
	u.PasswordHash, err = passhash.HashPassword(password, passhash.DefaultParams)
	return err
}

func (u *User) CheckPassword(password string) (bool, error) {
	return passhash.CheckPassword(password, u.PasswordHash)
}
