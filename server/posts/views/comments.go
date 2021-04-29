package views

import (
	"blogv2/posts/controllers"
	"blogv2/posts/models"
	"blogv2/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all comments related requests
type CommentsViews struct {
	DB *gorm.DB
}

var commentsController *controllers.CommentsController

// handle comments creatoin request and validate param
func (cv *CommentsViews) AddCommentView(c echo.Context) error {
	var comment models.Comment
	if err := c.Bind(&comment); err != nil {
		return echo.NewHTTPError(423, "unable to parse request bdy")
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid id")
	}

	userId := c.Request().Context().Value("user").(uint)
	comment.PostID = uint(postId)
	comment.UserID = userId

	httpErr := commentsController.AddComment(&comment, cv.DB)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(201, comment)
}

// retrieve comment by id and validate all url params
func (cv *CommentsViews) GetPostCommentsView(c echo.Context) error {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid id")
	}

	cursor := c.QueryParam("cursor")
	httpErr := utils.ValidateCursor(cursor)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid limit")
	}

	comments, hasMore := commentsController.GetPostComments(postId, limit, &cursor, cv.DB)

	return c.JSON(200, models.PaginatedComments{
		Comments: comments,
		HasMore:  hasMore,
	})
}

// handle delete comment request
func (cv *CommentsViews) DeleteCommentView(c echo.Context) error {
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(400, "invalid id")
	}

	userId := c.Request().Context().Value("user").(int)
	httpErr := commentsController.DeleteComment(commentId, userId, cv.DB)
	if err != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, "successfully deleted")
}
