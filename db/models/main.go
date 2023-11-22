package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id          uint   `gorm:"primaryKey"`
	Email       string `gorm:"unique"`
	DisplayName string
	Password    string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
