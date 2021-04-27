package views

import (
	"blogv2/posts/controllers"
	"blogv2/posts/models"
	"blogv2/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PostsViews struct {
	DB *gorm.DB
}

var postsController *controllers.PostsController

func (p *PostsViews) CreatePostView(c echo.Context) error {
	var post models.Post
	if err := c.Bind(&post); err != nil {
		return echo.NewHTTPError(423, "unable to parse request body")
	}

	userId, _ := strconv.Atoi(utils.UserIDFromToken(c))
	post.UserID = uint(userId)
	httpErr := postsController.CreatePost(&post, p.DB)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(201, post)
}

func (p *PostsViews) GetPostView(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := postsController.GetPost(id, p.DB)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, post)
}

func (p *PostsViews) DeletePostView(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	userId, _ := strconv.Atoi(utils.UserIDFromToken(c))
	err := postsController.DeletePost(id, userId, p.DB)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, "successfully deleted")
}
