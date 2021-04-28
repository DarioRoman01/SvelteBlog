package controllers

import (
	"blogv2/posts/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CommentsController struct{}

func (c *CommentsController) AddComment(comment *models.Comment, db *gorm.DB) *echo.HTTPError {
	if err := db.Create(&comment).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create post")
	}

	return nil
}

func (c *CommentsController) GetPostComments(postId, limit int, cursor *string, db *gorm.DB) ([]models.Comment, bool) {
	var comments []models.Comment
	if limit > 50 {
		limit = 50
	}
	limit++

	if cursor != nil {
		db.Raw(`
			SELECT c.*,
			FROM comments c
			WHERE post_id = ?
			AND c.created_at < ?
			ORDER BY p.created_at DESC
			LIMIT ?
		`, postId, cursor, limit).Find(&comments)
	} else {
		db.Raw(`
			SELECT p.*,
			FROM comments c
			WHERE post_id = ?
			ORDER BY p.created_at DESC
			LIMIT ?
		`, postId, limit).Find(&comments)
	}

	if len(comments) == 0 {
		return nil, false
	}
	if len(comments) == limit {
		return comments[0 : limit-1], true
	}

	return comments[0 : len(comments)-1], false
}

func (c *CommentsController) DeleteComment(id, userId int, db *gorm.DB) *echo.HTTPError {
	tx := db.Where("id = ? AND user_id = ?", id, userId).Delete(&models.Comment{})
	if tx.RowsAffected == 0 || tx.Error != nil {
		return echo.NewHTTPError(400, "post does not exists or you are not the owner")
	}

	return nil
}
