package views

import (
	"blogv2/users/controllers"
	"blogv2/users/models"
	"blogv2/utils"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UsersViews struct {
	DB *gorm.DB
}

var usersController *controllers.UsersController

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to read env: ", err)
	}
}

func (u *UsersViews) SignupView(c echo.Context) error {
	var userInput models.UserRegisterInput

	if err := c.Bind(&userInput); err != nil {
		return echo.NewHTTPError(423, "unable to parse request body")
	}

	isInvalid := utils.ValidateRegister(userInput)

	if isInvalid != nil {
		return c.JSON(isInvalid.Code, isInvalid.Message)
	}

	userCreated, err := usersController.CreateUser(&models.User{
		Email:       userInput.Email,
		PhoneNumber: userInput.PhoneNumber,
		Password:    userInput.Password,
	}, u.DB)

	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	send := utils.SendEmail(userCreated, "verify")
	if !send {
		return echo.NewHTTPError(500, "unable to send email")
	}

	return c.JSON(201, userCreated)
}

func (u *UsersViews) LoginView(c echo.Context) error {
	var userInput models.UserLoginInput

	if err := c.Bind(&userInput); err != nil {
		return echo.NewHTTPError(423, "unable to parse request body")
	}

	user, err := usersController.LoginByEmail(&userInput, u.DB)
	if err != nil {
		return c.JSON(err.Code, err.Message)
	}

	return c.JSON(200, user)
}

func (u *UsersViews) VerifyAccountView(c echo.Context) error {
	var body map[string]string
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return c.JSON(400, "unable to parse request body")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(body["token"], claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT-SECRET")), nil
	})
	if err != nil {
		return echo.NewHTTPError(423, "invalid token")
	}

	if claims["type"] != "verify" || !token.Valid {
		return echo.NewHTTPError(400, "invalid token")
	}

	err = u.DB.Model(&models.User{}).Where("id = ?", claims["user_id"]).Update("is_verfied = ?", true).Error
	if err != nil {
		return echo.NewHTTPError(500, "unable to update user status")
	}

	return c.JSON(200, "account verified succsefully")
}

func (u *UsersViews) ChangePasswordView(c echo.Context) error {
	var body map[string]string
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return c.JSON(423, "unable to parse request body")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(body["token"], claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT-SECRET")), nil
	})
	if err != nil {
		return echo.NewHTTPError(423, "invalid token")
	}

	if claims["type"] != "verify" || !token.Valid {
		return echo.NewHTTPError(400, "invalid token")

	} else if len(body["newPassword"]) < 8 {
		return echo.NewHTTPError(400, "password must be at least 8 characters")
	}

	user, httpErr := usersController.ChangePassword(claims["user_id"].(string), body["newPassword"], u.DB)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, user)
}

func (u *UsersViews) ForgotPasswordView(c echo.Context) error {
	var email map[string]string
	err := (&echo.DefaultBinder{}).BindBody(c, &email)
	if err != nil {
		return c.JSON(423, "unable to parse request body")
	}
	user := usersController.GetUserByEmail(email["email"], u.DB)
	if user == nil {
		return utils.InvalidInput("email", "email does not exist")
	}

	ok := utils.SendEmail(user, "change")

	if !ok {
		return echo.NewHTTPError(500, "unable to send email")
	}

	return c.JSON(200, "we send you and email please check your email")
}
