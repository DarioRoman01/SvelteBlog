package models

import (
	"time"

	"gorm.io/gorm"
)

// post model
type Post struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	UserID     uint           `json:"creatorId"`
	StateValue uint16         `json:"stateValue" gorm:"->"`
	Title      string         `json:"title"`
	Body       string         `json:"body"`
	Likes      uint           `json:"likes" gorm:"default=0"`
	Creator    string         `json:"creator" gorm:"->"`
	Comments   []Comment      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PostID;references:ID;"`
	Liked      []Like         `json:"-"`
}

// model for posts pagination
type PaginatedPosts struct {
	Posts   []Post `json:"posts"`
	HasMore bool   `json:"hasMore"`
}
