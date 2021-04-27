package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	UserID    uint           `json:"userID"`
	Username  string         `json:"username" gorm:"unique"`
	Biography string         `json:"biography"`
	Followers uint           `json:"followers" gorm:"default=0"`
	Posted    uint           `json:"posted" gorm:"default=0"`
}
