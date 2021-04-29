package controllers

import (
	"blogv2/posts/models"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all posts related functions in the db
type PostsController struct{}

// store the post in the db
func (p *PostsController) CreatePost(post *models.Post, db *gorm.DB) *echo.HTTPError {
	if err := db.Create(&post).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create post.")
	}

	return nil
}

// retrieve post by id
func (p *PostsController) GetPost(id int, db *gorm.DB) (*models.Post, *echo.HTTPError) {
	var post models.Post
	db.Model(&models.Post{}).Where("id = ?", id).Find(&post)

	if post.ID == 0 {
		return nil, echo.NewHTTPError(404, "post does not exist")
	}

	return &post, nil
}

// delete post from the db
func (p *PostsController) DeletePost(id int, userID int, db *gorm.DB) *echo.HTTPError {
	tx := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Post{})
	if tx.RowsAffected == 0 || tx.Error != nil {
		return echo.NewHTTPError(400, "post does not exist or you are not the owner")
	}

	return nil
}

// retrieve all post and paginate them
func (p *PostsController) GetPosts(limit int, cursor *string, userId int, db *gorm.DB) ([]models.Post, bool) {
	var posts []models.Post
	if limit > 50 {
		limit = 50
	}
	limit++

	if cursor != nil {
		db.Raw(`
			SELECT p.*,
			( SELECT "value" from "likes" 
			WHERE "user_id" = ? and "post_id" = p.id) as "StateValue"
			FROM posts p
			WHERE p.created_at < ?
			ORDER BY p.created_at DESC
			LIMIT ?
		`, userId, cursor, limit).Find(&posts)
	} else {
		db.Raw(`
			SELECT p.*,
			( SELECT "value" from "updoots" 
			WHERE "user_id" = ? and "post_id" = p.id) as "StateValue"
			FROM posts p
			ORDER BY p.created_at DESC
			LIMIT ?
		`, userId, limit).Find(&posts)
	}
	if len(posts) == 0 {
		return nil, false
	}
	if len(posts) == limit {
		return posts[0 : limit-1], true
	}

	return posts[0 : len(posts)-1], false
}

// set user like or quit their like depending if he liked the post before
func (p *PostsController) SetLike(postId, userId, value int, db *gorm.DB) bool {
	var like models.Like
	isLike := value != -1
	var realValue uint16

	if isLike {
		realValue = 1
	} else {
		realValue = 0
	}

	db.Table("likes").Where("user_id = ? and post_id = ?", userId, postId).Find(&like)

	// user is liked the post before and
	// they are changing their like
	if like.PostID != 0 && like.Value != realValue {
		query := fmt.Sprintf(`
			START TRANSACTION;

			UPDATE "likes"
			SET value = %d
			WHERE post_id = %d AND user_id = %d; 

			UPDATE "posts"
			SET Likes = Likes - 1
			WHERE posts.id = %d;
			
			COMMIT;
		`, realValue, postId, userId, postId)

		if err := db.Exec(query).Error; err != nil {
			return false
		}

		// user has never liked before
	} else if like.PostID == 0 {
		query := fmt.Sprintf(`
			START TRANSACTION;

			INSERT INTO "likes" ("user_id", "post_id", "value")
			values(%d, %d, %d);

			UPDATE "posts"
			SET Likes = Likes + %d
			WHERE posts.id = %d;

			COMMIT;
		`, userId, postId, realValue, realValue, postId)

		if err := db.Exec(query).Error; err != nil {
			return false
		}
	}

	return true
}

// retrieve the posts of the given user paginated
func (p *PostsController) GetUserPosts(limit, userId, profileId int, cursor *string, db *gorm.DB) ([]models.Post, bool) {
	var posts []models.Post

	if limit > 50 {
		limit = 50
	}
	limit++

	if cursor != nil {
		db.Raw(`
			SELECT p.*,
			( SELECT "value" from "likes" 
			WHERE "user_id" = ? and "post_id" = p.id) as "StateValue"
			FROM posts p
			WHERE user_id = ?
			AND p.created_at < ?
			ORDER BY p.created_at DESC
			LIMIT ?
		`, userId, profileId, cursor, limit).Find(&posts)
	} else {
		db.Raw(`
			SELECT p.*,
			( SELECT "value" from "likes" 
			WHERE "user_id" = ? and "post_id" = p.id) as "StateValue"
			FROM posts p
			WHERE user_id = ?
			ORDER BY p.created_at DESC
			LIMIT ?
		`, userId, profileId, limit).Find(&posts)
	}
	if len(posts) == 0 {
		return nil, false
	}
	if len(posts) == limit {
		return posts[0 : limit-1], true
	}

	return posts[0 : len(posts)-1], false
}
