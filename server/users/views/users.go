package views

import (
	"blogv2/users/controllers"
	"blogv2/users/models"
	"blogv2/utils"
	"log"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// gruops all users views
type UsersViews struct {
	DB *gorm.DB
}

// shortcut for not writing controllers.UsersControllers
var usersController *controllers.UsersController

// read env to let all functions get .env variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to read env: ", err)
	}
}

// hande signup requests and validate all fields to then call the controller
func (u *UsersViews) SignupView(c echo.Context) error {
	var userInput models.UserRegisterInput
	if err := c.Bind(&userInput); err != nil {
		return utils.RequestBodyError
	}

	httpErr := utils.ValidateRegister(userInput)
	if httpErr != nil {
		return httpErr
	}

	userCreated, httpErr := usersController.CreateUser(&models.User{
		Email:       userInput.Email,
		PhoneNumber: userInput.PhoneNumber,
		Password:    userInput.Password,
	}, u.DB)
	if httpErr != nil {
		return httpErr
	}

	send := utils.SendVerificationEmail(userCreated)
	if !send {
		return echo.NewHTTPError(500, "unable to send email")
	}

	return c.JSON(201, userCreated)
}

// hande login requests and generate a session for the user with jwt token
func (u *UsersViews) LoginView(c echo.Context) error {
	var userInput models.UserLoginInput

	if err := c.Bind(&userInput); err != nil {
		return utils.RequestBodyError
	}

	user, httpErr := usersController.LoginByEmail(&userInput, u.DB)
	if httpErr != nil {
		return httpErr
	}

	token, httpErr := utils.GenerateToken(user.ID, utils.SESSION)
	if httpErr != nil {
		return httpErr
	}

	session, _ := utils.Store.Get(c.Request(), "jwt")
	session.Values["token"] = token
	session.Save(c.Request(), c.Response().Writer)
	return c.JSON(200, user)
}

// handle users verifycation validating the token send by the user
// and generate a cookie with a new jwt token for session
func (u *UsersViews) VerifyAccountView(c echo.Context) error {
	var body map[string]string
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return utils.RequestBodyError
	}

	_, exists := body["token"]
	if !exists {
		return utils.InvalidInput("token", "token field is required")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(body["token"], claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT-SECRET")), nil
	})
	if err != nil {
		return echo.NewHTTPError(400, "invalid token")
	}

	if claims["type"] != utils.VERIFY || !token.Valid {
		return echo.NewHTTPError(400, "invalid token")
	}

	err = u.DB.Table("users").Where("id = ?", int(claims["user_id"].(float64))).Update("is_verified", true).Error
	if err != nil {
		return echo.NewHTTPError(500, "unable to update user status")
	}

	sessionToken, httpErr := utils.GenerateToken(uint(claims["user_id"].(float64)), utils.SESSION)
	if httpErr != nil {
		return httpErr
	}

	session, _ := utils.Store.Get(c.Request(), "jwt")
	session.Values["token"] = sessionToken
	session.Save(c.Request(), c.Response().Writer)
	return c.JSON(200, "account verified succsefully")
}

// handle change password requests for users
// validate jwt token send by the user
func (u *UsersViews) ChangePasswordView(c echo.Context) error {
	var body map[string]string
	err := (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		return utils.RequestBodyError
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(body["token"], claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT-SECRET")), nil
	})
	if err != nil {
		return echo.NewHTTPError(400, "invalid token")
	}

	if claims["type"] != utils.CHANGE || !token.Valid {
		return echo.NewHTTPError(400, "invalid token")

	} else if len(body["newPassword"]) < 8 {
		return echo.NewHTTPError(400, "password must be at least 8 characters")
	}

	user, httpErr := usersController.ChangePassword(claims["user_id"].(string), body["newPassword"], u.DB)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, user)
}

// handle forgot password request and validate the the
// introduced email is use in the app
func (u *UsersViews) ForgotPasswordView(c echo.Context) error {
	var email map[string]string
	err := (&echo.DefaultBinder{}).BindBody(c, &email)
	if err != nil {
		return utils.RequestBodyError
	}

	user := usersController.GetUserByEmail(email["email"], u.DB)
	if user == nil {
		return utils.InvalidInput("email", "email does not exist")
	}

	ok := utils.SendChangePasswordEmail(user)
	if !ok {
		return echo.NewHTTPError(500, "unable to send email")
	}

	return c.JSON(200, "we send you and email please check your email")
}

// delete the user session
func (u *UsersViews) LogoutView(c echo.Context) error {
	session, _ := utils.Store.Get(c.Request(), "jwt")
	session.Options.MaxAge = -1
	err := session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return echo.NewHTTPError(500, "unable to delete session")
	}

	return c.JSON(200, "succesfully logout")
}

// retrieve all followers of the user paginated
// and validate params and query params
func (u *UsersViews) UserFollowersViews(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return utils.LimitError
	}

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	cursor := c.QueryParam("cursor")
	httpErr := utils.ValidateCursor(cursor)
	if httpErr != nil {
		return httpErr
	}

	users, hasMore := usersController.GetUserFollowers(userId, limit, cursor, u.DB)
	return c.JSON(200, models.PaginatedUsers{Users: users, HasMore: hasMore})
}
