package views

import (
	"blogv2/posts/controllers"
	"blogv2/posts/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CommentsViews struct {
	DB *gorm.DB
}

var commentsController *controllers.CommentsController

func (cv *CommentsViews) AddCommentView(c echo.Context) error {
	var comment models.Comment
	if err := c.Bind(&comment); err != nil {
		return echo.NewHTTPError(423, "unable to parse request bdy")
	}

	postId, _ := strconv.Atoi(c.Param("id"))
	userId := c.Request().Context().Value("user").(uint)
	comment.PostID = uint(postId)
	comment.UserID = userId

	httpErr := commentsController.AddComment(&comment, cv.DB)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(201, comment)
}

func (cv *CommentsViews) GetPostCommentsView(c echo.Context) error {
	postId, _ := strconv.Atoi(c.Param("id"))
	cursor := c.QueryParam("cursor")
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

func (cv *CommentsViews) DeleteCommentView(c echo.Context) error {
	commentId, _ := strconv.Atoi(c.Param("id"))
	userId := c.Request().Context().Value("user").(int)
	err := commentsController.DeleteComment(commentId, userId, cv.DB)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, "successfully deleted")
}
