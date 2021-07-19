package routes

import (
	"blogv2/db"
	pControllers "blogv2/posts/controllers"
	pViews "blogv2/posts/views"
	uControllers "blogv2/users/controllers"
	"blogv2/users/views"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetRoutes(e *echo.Echo) {
	db, err := db.Connect()
	if err != nil {
		log.Fatal("unable to connect to postgres: ", err)
	}

	// handlers
	usersViews := views.NewUsersViews(uControllers.NewUserController(db))
	postsViews := pViews.NewPostsViews(pControllers.NewPostController(db))

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
	e.POST("/logout", usersViews.LogoutView)

	// profile views
	e.GET("/me", usersViews.MeView)
	e.POST("/profile", usersViews.CreateProfileView)
	e.POST("/profile/:id/follow", usersViews.FollowView)
	e.PATCH("/profile/:id", usersViews.UpdateProfileView)
	e.GET("/profile/:username", usersViews.GetProfileView)
	e.GET("/profile/:id/posts", postsViews.GetUserPostsView)

	// posts views
	e.GET("/posts", postsViews.GetPostsView)
	e.POST("/posts", postsViews.CreatePostView)
	e.GET("/posts/:id", postsViews.GetPostView)
	e.DELETE("/posts/:id", postsViews.DeletePostView)
	e.POST("/posts/:id/like", postsViews.ToggleLikeView)
	e.PATCH("/posts/:id", postsViews.UpdatePostView)
	e.POST("/posts/:id/comments", postsViews.AddCommentView)
	e.GET("/posts/:id/comments", postsViews.GetPostCommentsView)
	e.DELETE("/posts/:id/comments", postsViews.DeleteCommentView)

}
