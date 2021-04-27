package controllers

import (
	"blogv2/posts/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PostsController struct{}

func (p *PostsController) CreatePost(post *models.Post, db *gorm.DB) *echo.HTTPError {
	if err := db.Create(&post).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create post.")
	}

	return nil
}

func (p *PostsController) GetPost(id int, db *gorm.DB) (*models.Post, *echo.HTTPError) {
	var post models.Post
	db.Model(&models.Post{}).Where("id = ?", id).Find(&post)

	if post.ID == 0 {
		return nil, echo.NewHTTPError(404, "post does not exist")
	}

	return &post, nil
}

func (p *PostsController) DeletePost(id int, userID int, db *gorm.DB) *echo.HTTPError {
	tx := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Post{})
	if tx.RowsAffected == 0 || tx.Error != nil {
		return echo.NewHTTPError(400, "post does not exist or you are not the owner")
	}

	return nil
}
