package models

import (
	"time"

	"gorm.io/gorm"
)

// profile model holds all user public info
type Profile struct {
	UserID      uint           `json:"userID" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Username    string         `json:"username" gorm:"unique"`
	Biography   string         `json:"biography"`
	Followers   uint           `json:"followers" gorm:"default=0"`
	Posted      uint           `json:"posted" gorm:"default=0"`
	FollowState bool           `json:"followState" gorm:"->"`
}
