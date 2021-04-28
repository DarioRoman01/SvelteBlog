package models

import "time"

type Comment struct {
	UserID    uint      `json:"creatorId" gorm:"primaryKey"`
	PostID    uint      `json:"postId" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
}

type PaginatedComments struct {
	Comments []Comment `json:"comments"`
	HasMore  bool      `json:"hasMore"`
}
