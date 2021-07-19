package views

import (
	"blogv2/posts/controllers"
	"blogv2/posts/models"
	"blogv2/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

// group all requests related to posts
type PostsViews struct {
	controller *controllers.PostsController
}

func NewPostsViews(controller *controllers.PostsController) *PostsViews {
	return &PostsViews{controller}
}

// handle create post request
func (p *PostsViews) CreatePostView(c echo.Context) error {
	var post models.Post
	if err := c.Bind(&post); err != nil {
		return utils.RequestBodyError
	}

	userId := c.Request().Context().Value("user").(uint)
	post.UserID = userId
	httpErr := p.controller.CreatePost(&post)
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
	post, httpErr := p.controller.GetPost(id, userId)
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
	httpErr := p.controller.DeletePost(id, int(userId))
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
	posts, hasMore := p.controller.GetPosts(limit, cursor, int(userId))
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
	likedOrDisliked := p.controller.SetLike(postId, int(userId), body["value"])

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

	posts, hasMore := p.controller.GetUserPosts(limit, int(userId), profileId, cursor)
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
	post, httpErr := p.controller.UpdatePost(postId, userId, data)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, post)
}

// handle comments creatoin request and validate param
func (p *PostsViews) AddCommentView(c echo.Context) error {
	var comment models.Comment
	if err := c.Bind(&comment); err != nil {
		return utils.RequestBodyError
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	userId := c.Request().Context().Value("user").(uint)
	comment.PostID = uint(postId)
	comment.UserID = userId

	httpErr := p.controller.AddComment(&comment)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(201, comment)
}

// retrieve comments by post id and validate all url params
func (p *PostsViews) GetPostCommentsView(c echo.Context) error {
	postId, err := strconv.Atoi(c.Param("id"))
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

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return utils.LimitError
	}

	comments, hasMore := p.controller.GetPostComments(postId, limit, cursor)
	return c.JSON(200, models.PaginatedComments{
		Comments: comments,
		HasMore:  hasMore,
	})
}

// handle delete comment request
func (p *PostsViews) DeleteCommentView(c echo.Context) error {
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	userId := c.Request().Context().Value("user").(int)
	httpErr := p.controller.DeleteComment(commentId, userId)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, "successfully deleted")
}
