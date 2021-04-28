package views

import (
	"blogv2/posts/controllers"
	"blogv2/posts/models"
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

	userId := c.Request().Context().Value("user").(uint)
	post.UserID = userId
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
	userId := c.Request().Context().Value("user").(int)
	err := postsController.DeletePost(id, userId, p.DB)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, "successfully deleted")
}

func (p *PostsViews) GetPostsView(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid limit")
	}

	cursor := c.QueryParam("cursor")
	userId := c.Request().Context().Value("user").(int)
	posts, hasMore := postsController.GetPosts(limit, &cursor, userId, p.DB)

	return c.JSON(200, models.PaginatedPosts{Posts: posts, HasMore: hasMore})
}

func (p *PostsViews) ToggleLikeView(c echo.Context) error {
	var body map[string]int
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return c.JSON(423, "unable to parse request body")
	}

	postId, _ := strconv.Atoi(c.Param("id"))
	userId := c.Request().Context().Value("user").(int)
	likedOrDisliked := postsController.SetLike(postId, userId, body["value"], p.DB)

	if !likedOrDisliked {
		return echo.NewHTTPError(500, "unable to set like")
	}

	return c.JSON(201, "liked successfully")
}

func (p *PostsViews) GetUserPostsView(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid limit")
	}

	userId := c.Request().Context().Value("user").(int)
	profileId, _ := strconv.Atoi(c.Param("id"))
	cursor := c.QueryParam("cursor")

	posts, hasMore := postsController.GetUserPosts(limit, userId, profileId, &cursor, p.DB)

	return c.JSON(200, models.PaginatedPosts{Posts: posts, HasMore: hasMore})
}
