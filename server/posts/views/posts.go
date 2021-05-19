package views

import (
	"blogv2/posts/controllers"
	"blogv2/posts/models"
	"blogv2/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all requests related to posts
type PostsViews struct {
	DB *gorm.DB
}

var postsController *controllers.PostsController

// handle create post request
func (p *PostsViews) CreatePostView(c echo.Context) error {
	var post models.Post
	if err := c.Bind(&post); err != nil {
		return utils.RequestBodyError
	}

	userId := c.Request().Context().Value("user").(uint)
	post.UserID = userId
	httpErr := postsController.CreatePost(&post, p.DB)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(201, post)
}

// retrieve post by id and validate id in url params
func (p *PostsViews) GetPostView(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(423, "invalid id")
	}
	userId := c.Request().Context().Value("user").(uint)
	post, httpErr := postsController.GetPost(id, userId, p.DB)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, post)
}

// handle delete post request and validate url params
func (p *PostsViews) DeletePostView(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	userId := c.Request().Context().Value("user").(uint)
	httpErr := postsController.DeletePost(id, int(userId), p.DB)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, "successfully deleted")
}

// handle list posts validating url query params for pagination
func (p *PostsViews) GetPostsView(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return utils.LimitError
	}

	cursor := c.QueryParam("cursor")
	if cursor != "" {
		httpErr := utils.ValidateCursor(cursor)
		if httpErr != nil {
			return httpErr
		}
	}

	userId := c.Request().Context().Value("user").(uint)
	posts, hasMore := postsController.GetPosts(limit, cursor, int(userId), p.DB)

	return c.JSON(200, models.PaginatedPosts{Posts: posts, HasMore: hasMore})
}

// handle like request and validta url params
func (p *PostsViews) ToggleLikeView(c echo.Context) error {
	var body map[string]int
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return utils.RequestBodyError
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	userId := c.Request().Context().Value("user").(uint)
	likedOrDisliked := postsController.SetLike(postId, int(userId), body["value"], p.DB)

	if !likedOrDisliked {
		return echo.NewHTTPError(500, "unable to set like")
	}

	return c.JSON(201, "liked successfully")
}

// retrieve all users posts with pagination
func (p *PostsViews) GetUserPostsView(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid limit")
	}

	userId := c.Request().Context().Value("user").(uint)
	profileId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	cursor := c.QueryParam("cursor")
	if cursor != "" {
		httpErr := utils.ValidateCursor(cursor)
		if httpErr != nil {
			return httpErr
		}
	}

	posts, hasMore := postsController.GetUserPosts(limit, int(userId), profileId, cursor, p.DB)
	return c.JSON(200, models.PaginatedPosts{Posts: posts, HasMore: hasMore})
}

func (p *PostsViews) UpdatePostView(c echo.Context) error {
	var data models.Post
	if err := c.Bind(&data); err != nil {
		return utils.RequestBodyError
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	userId := c.Request().Context().Value("user").(uint)
	post, httpErr := postsController.UpdatePost(postId, userId, data, p.DB)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, post)
}
