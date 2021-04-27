package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	UserID     uint           `json:"creator_id"`
	StateValue uint16         `json:"stateValue" gorm:"-:migration"`
	Title      string         `json:"title"`
	Body       string         `json:"body"`
	Likes      uint           `json:"likes" gorm:"default=0"`
	Comments   []Comment      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PostID;references:ID;"`
	Liked      []Like
}
