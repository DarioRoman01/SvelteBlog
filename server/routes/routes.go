package routes

import (
	"blogv2/db"
	pModels "blogv2/posts/models"
	pViews "blogv2/posts/views"
	"blogv2/users/models"
	"blogv2/users/views"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetRoutes(e *echo.Echo) {
	psql, err := db.Connect()
	if err != nil {
		log.Fatal("unable to connect to postgres: ", err)
	}

	psql.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&pModels.Post{},
		&pModels.Comment{},
		&pModels.Like{},
	)

	// handlers
	usersViews := &views.UsersViews{DB: psql}
	profileViews := &views.ProfileViews{DB: psql}
	postsViews := &pViews.PostsViews{DB: psql}
	commentViews := &pViews.CommentsViews{DB: psql}

	// middlewares
	e.Use(CORSconfig())
	e.Use(IsAuth)
	e.Use(middleware.RemoveTrailingSlash())

	// user views
	e.POST("/register", usersViews.SignupView)
	e.POST("/login", usersViews.LoginView)
	e.POST("/change-password", usersViews.ChangePasswordView)
	e.POST("/forgot-password", usersViews.ForgotPasswordView)
	e.POST("/verify", usersViews.VerifyAccountView)

	// profile views
	e.PATCH("/profile/:id", profileViews.UpdateProfileView)
	e.POST("/profile", profileViews.CreateProfileView)
	e.GET("/profile/:username", profileViews.GetProfileView)
	e.GET("/profile/:id/posts", postsViews.GetUserPostsView)
	e.POST("/profile/:id/follow", profileViews.FollowView)

	// posts views
	e.GET("/posts", postsViews.GetPostsView)
	e.POST("/posts", postsViews.CreatePostView)
	e.GET("/posts/:id", postsViews.GetPostView)
	e.DELETE("/posts/:id", postsViews.DeletePostView)
	e.POST("/posts/:id/like", postsViews.ToggleLikeView)
	e.POST("/posts/:id/comments", commentViews.AddCommentView)
	e.GET("/posts/:id/comments", commentViews.GetPostCommentsView)
	e.DELETE("/posts/:id/comments", commentViews.DeleteCommentView)

}
