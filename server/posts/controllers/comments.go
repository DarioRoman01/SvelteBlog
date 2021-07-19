package controllers

import (
	"blogv2/posts/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all functions related on comments on the db
type CommentsController struct {
	db *gorm.DB
}

// create comment in the db
func (c *CommentsController) AddComment(comment *models.Comment) *echo.HTTPError {
	if err := c.db.Table("comments").Create(&comment).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create post")
	}

	return nil
}

// retrieve all posts comments and paginate them
func (c *CommentsController) GetPostComments(postId, limit int, cursor string) ([]models.Comment, bool) {
	var comments []models.Comment
	if limit > 50 {
		limit = 50
	}
	limit++

	if cursor != "" {
		c.db.Raw(`
			SELECT c.*,
			(SELECT username FROM "profiles"
			WHERE user_id = c.user_id) as "Creator"
			FROM comments c
			WHERE c.post_id = ?
			AND c.created_at < ?
			ORDER BY c.created_at DESC
			LIMIT ?
		`, postId, cursor, limit).Find(&comments)
	} else {
		c.db.Raw(`
			SELECT c.*,
			(SELECT username FROM "profiles"
			WHERE user_id = c.user_id) as "Creator"
			FROM comments c
			WHERE c.post_id = ?
			ORDER BY c.created_at DESC
			LIMIT ?
		`, postId, limit).Find(&comments)
	}

	if len(comments) == 0 {
		return nil, false
	}
	if len(comments) == limit {
		return comments[0 : limit-1], true
	}

	return comments, false
}

// delete comment from the db
func (c *CommentsController) DeleteComment(id, userId int) *echo.HTTPError {
	tx := c.db.Where("id = ? AND user_id = ?", id, userId).Delete(&models.Comment{})
	if tx.RowsAffected == 0 || tx.Error != nil {
		return echo.NewHTTPError(400, "post does not exists or you are not the owner")
	}

	return nil
}
