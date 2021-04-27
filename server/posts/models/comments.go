package models

import "time"

type Comment struct {
	UserID    uint      `json:"creator_id" gorm:"primaryKey"`
	PostID    uint      `json:"post" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
}
