package models

import "time"

// Comment model
type Comment struct {
	UserID    uint      `json:"creatorId" gorm:"primaryKey"`
	PostID    uint      `json:"postId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Creator   string    `json:"creator"`
}

// model for comments pagination
type PaginatedComments struct {
	Comments []Comment `json:"comments"`
	HasMore  bool      `json:"hasMore"`
}
