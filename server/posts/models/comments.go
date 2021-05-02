package models

import "time"

// Comment model
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"creatorId"`
	PostID    uint      `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Creator   string    `json:"creator" gorm:"->"`
}

// model for comments pagination
type PaginatedComments struct {
	Comments []Comment `json:"comments"`
	HasMore  bool      `json:"hasMore"`
}
