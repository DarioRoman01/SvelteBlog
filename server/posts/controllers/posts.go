package controllers

import (
	"blogv2/posts/models"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all posts related functions in the db
type PostsController struct {
	db *gorm.DB
}

func NewPostController(db *gorm.DB) *PostsController {
	return &PostsController{db}
}

// store the post in the db and update profile posted counter with transaction
func (p *PostsController) CreatePost(post *models.Post) *echo.HTTPError {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return echo.NewHTTPError(500, "unable to start transaction")
	}

	if err := tx.Table("posts").Create(&post).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(500, "unable to create post.")
	}

	err := tx.Exec(`UPDATE "profiles" SET posted = posted + 1 WHERE user_id = ?`, post.UserID).Error
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(500, "unable to update profile")
	}
	if err := tx.Commit().Error; err != nil {
		return echo.NewHTTPError(500, "unable to ccommit")
	}

	return nil
}

// retrieve post by id
func (p *PostsController) GetPost(id int, userId uint) (*models.Post, *echo.HTTPError) {
	var post models.Post
	p.db.Raw(`
		SELECT p.*,
		(SELECT "value" from "likes" 
		WHERE "user_id" = ? and "post_id" = p.id) as "StateValue",
		(SELECT username FROM "profiles"
		WHERE user_id = p.user_id) as "Creator"
		FROM posts p
		WHERE p.id = ?
	`, userId, id).Find(&post)

	if post.ID == 0 {
		return nil, echo.NewHTTPError(404, "post does not exist")
	}

	return &post, nil
}

// delete post from the db with a transaction to handle user published posts count
func (p *PostsController) DeletePost(id int, userID int) *echo.HTTPError {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return echo.NewHTTPError(500, "unable to start transaction")
	}

	result := tx.Exec(`DELETE FROM "posts" WHERE id = ? AND user_id = ?`, id, userID)
	if result.RowsAffected == 0 || result.Error != nil {
		tx.Rollback()
		return echo.NewHTTPError(400, "post does not exists or you are not the owner")
	}

	err := tx.Exec(`UPDATE "profiles" SET posted = posted - 1 WHERE user_id = ?`, userID).Error
	if err != nil {
		tx.Rollback()
		return echo.NewHTTPError(500, "unable to update profile")
	}

	if err := tx.Commit().Error; err != nil {
		return echo.NewHTTPError(500, "unable to ccommit")
	}

	return nil
}

// retrieve all post and paginate them
func (p *PostsController) GetPosts(limit int, cursor string, userId int) ([]models.Post, bool) {
	var posts []models.Post
	if limit > 50 {
		limit = 50
	}
	limit++

	if cursor != "" {
		p.db.Raw(`
			SELECT p.*,
			(SELECT "value" from "likes" 
			WHERE "user_id" = ? and "post_id" = p.id) as "StateValue",
			(SELECT username FROM "profiles"
			WHERE user_id = p.user_id) as "Creator"
			FROM posts p
			WHERE p.created_at < ?
			ORDER BY p.created_at DESC
			LIMIT ?
		`, userId, cursor, limit).Find(&posts)
	} else {
		p.db.Raw(`
			SELECT p.*,
			( SELECT "value" from "likes" 
			WHERE "user_id" = ? and "post_id" = p.id) as "StateValue",
			(SELECT "username" FROM "profiles"
			WHERE user_id = p.user_id) as "Creator"
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

	return posts, false
}

// set user like or quit their like depending if he liked the post before
func (p *PostsController) SetLike(postId, userId, value int) bool {
	var like models.Like
	isLike := value != -1
	var realValue uint16

	if isLike {
		realValue = 1
	} else {
		realValue = 0
	}

	p.db.Table("likes").Where("user_id = ? and post_id = ?", userId, postId).Find(&like)

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

		if err := p.db.Exec(query).Error; err != nil {
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

		if err := p.db.Exec(query).Error; err != nil {
			return false
		}
	}

	return true
}

// retrieve the posts of the given user paginated
func (p *PostsController) GetUserPosts(limit, userId, profileId int, cursor string) ([]models.Post, bool) {
	var posts []models.Post

	if limit > 50 {
		limit = 50
	}
	limit++

	if cursor != "" {
		p.db.Raw(`
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
		p.db.Raw(`
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

	return posts, false
}

func (p *PostsController) UpdatePost(postID int, userID uint, data models.Post) (*models.Post, *echo.HTTPError) {
	var post models.Post
	p.db.Table("posts").Where("id = ?", postID).Find(&post)
	if post.ID == 0 {
		return nil, echo.NewHTTPError(404, "post not found")
	}

	if post.UserID != userID {
		return nil, echo.NewHTTPError(403, "you are not allowed to perform this action")
	}

	if err := p.db.Model(&post).Updates(data).Error; err != nil {
		return nil, echo.NewHTTPError(500, "unable to update post :(")
	}

	return &post, nil
}

func (p *PostsController) AddComment(comment *models.Comment) *echo.HTTPError {
	if err := p.db.Table("comments").Create(&comment).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create post")
	}

	return nil
}

// retrieve all posts comments and paginate them
func (p *PostsController) GetPostComments(postId, limit int, cursor string) ([]models.Comment, bool) {
	var comments []models.Comment
	if limit > 50 {
		limit = 50
	}
	limit++

	if cursor != "" {
		p.db.Raw(`
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
		p.db.Raw(`
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
func (p *PostsController) DeleteComment(id, userId int) *echo.HTTPError {
	tx := p.db.Where("id = ? AND user_id = ?", id, userId).Delete(&models.Comment{})
	if tx.RowsAffected == 0 || tx.Error != nil {
		return echo.NewHTTPError(400, "post does not exists or you are not the owner")
	}

	return nil
}
