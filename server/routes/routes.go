package routes

import (
	"blogv2/db"
	"blogv2/users/models"
	"blogv2/users/views"
	"log"

	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	psql, err := db.Connect()
	if err != nil {
		log.Fatal("unable to connect to postgres: ", err)
	}

	psql.AutoMigrate(&models.User{}, &models.Profile{})
	usersViews := &views.UsersViews{DB: psql}

	e.POST("/register", usersViews.SignupView)
	e.POST("/login", usersViews.LoginView)
	e.POST("/change-password", usersViews.ChangePasswordView)
	e.POST("/forgot-password", usersViews.ForgotPasswordView)
	e.POST("/verify", usersViews.VerifyAccountView)
}
