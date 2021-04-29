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
	err := c.Bind(&post)
	utils.CheckRequestBodyError(err)

	userId := c.Request().Context().Value("user").(uint)
	post.UserID = userId
	httpErr := postsController.CreatePost(&post, p.DB)
	utils.CheckHttpError(httpErr)

	return c.JSON(201, post)
}

// retrieve post by id and validate id in url params
func (p *PostsViews) GetPostView(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	utils.CheckIDParamError(err)

	post, httpErr := postsController.GetPost(id, p.DB)
	utils.CheckHttpError(httpErr)

	return c.JSON(200, post)
}

// handle delete post request and validate url params
func (p *PostsViews) DeletePostView(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	utils.CheckIDParamError(err)
	userId := c.Request().Context().Value("user").(int)
	httpErr := postsController.DeletePost(id, userId, p.DB)
	utils.CheckHttpError(httpErr)

	return c.JSON(200, "successfully deleted")
}

// handle list posts validating url query params for pagination
func (p *PostsViews) GetPostsView(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	utils.CheckLimitParamError(err)

	cursor := c.QueryParam("cursor")
	utils.ValidateCursor(cursor)

	userId := c.Request().Context().Value("user").(int)
	posts, hasMore := postsController.GetPosts(limit, &cursor, userId, p.DB)

	return c.JSON(200, models.PaginatedPosts{Posts: posts, HasMore: hasMore})
}

// handle like request and validta url params
func (p *PostsViews) ToggleLikeView(c echo.Context) error {
	var body map[string]int
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	utils.CheckRequestBodyError(err)

	postId, err := strconv.Atoi(c.Param("id"))
	utils.CheckIDParamError(err)
	userId := c.Request().Context().Value("user").(int)
	likedOrDisliked := postsController.SetLike(postId, userId, body["value"], p.DB)

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

	userId := c.Request().Context().Value("user").(int)
	profileId, err := strconv.Atoi(c.Param("id"))
	utils.CheckIDParamError(err)

	cursor := c.QueryParam("cursor")
	utils.ValidateCursor(cursor)

	posts, hasMore := postsController.GetUserPosts(limit, userId, profileId, &cursor, p.DB)
	return c.JSON(200, models.PaginatedPosts{Posts: posts, HasMore: hasMore})
}
