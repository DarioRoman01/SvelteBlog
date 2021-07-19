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
)

// gruops all users views
type UsersViews struct {
	controller *controllers.UsersController
}

func NewUsersViews(controller *controllers.UsersController) *UsersViews {
	return &UsersViews{controller}
}

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

	user := &models.User{
		Email:       userInput.Email,
		PhoneNumber: userInput.PhoneNumber,
		Password:    userInput.Password,
	}

	httpErr = u.controller.CreateUser(user)
	if httpErr != nil {
		return httpErr
	}

	send := utils.SendVerificationEmail(user)
	if !send {
		return echo.NewHTTPError(500, "unable to send email")
	}

	return c.JSON(201, user)
}

// hande login requests and generate a session for the user with jwt token
func (u *UsersViews) LoginView(c echo.Context) error {
	var userInput models.UserLoginInput

	if err := c.Bind(&userInput); err != nil {
		return utils.RequestBodyError
	}

	user, httpErr := u.controller.LoginByEmail(&userInput)
	if httpErr != nil {
		return httpErr
	}

	token, httpErr := utils.GenerateToken(user.ID, utils.SESSION)
	if httpErr != nil {
		return httpErr
	}

	utils.SetToken(c, token)
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
	id := int(claims["user_id"].(float64))

	httpErr := u.controller.UpdateUserStatus(id)
	if httpErr != nil {
		return echo.NewHTTPError(500, "unable to update user status")
	}

	sessionToken, httpErr := utils.GenerateToken(uint(id), utils.SESSION)
	if httpErr != nil {
		return httpErr
	}

	utils.SetToken(c, sessionToken)
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

	if claims["type"] != utils.CHANGE || !token.Valid {
		return echo.NewHTTPError(400, "invalid token")

	} else if len(body["newPassword"]) < 8 {
		return echo.NewHTTPError(400, "password must be at least 8 characters")
	}

	user, httpErr := u.controller.ChangePassword(claims["user_id"].(string), body["newPassword"])
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

	user := u.controller.GetUserByEmail(email["email"])
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
	if err := utils.Unabletoken(c); err != nil {
		return err
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

	users, hasMore := u.controller.GetUserFollowers(userId, limit, cursor)
	return c.JSON(200, models.PaginatedUsers{Users: users, HasMore: hasMore})
}

func (u *UsersViews) CreateProfileView(c echo.Context) error {
	var profile models.Profile
	if err := c.Bind(&profile); err != nil {
		return utils.RequestBodyError
	}

	userId := c.Request().Context().Value("user").(uint)
	profile.UserID = userId
	httpErr := u.controller.CreateProfile(&profile)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(201, profile)
}

// handle profile updates
func (u *UsersViews) UpdateProfileView(c echo.Context) error {
	var profile models.Profile
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	if err = c.Bind(&profile); err != nil {
		return utils.RequestBodyError
	}

	userId := c.Request().Context().Value("user").(uint)
	newProfile, httpErr := u.controller.UpdateProfile(userId, uint(id), &profile)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, newProfile)
}

// retrieve profile with given username
func (u *UsersViews) GetProfileView(c echo.Context) error {
	username := c.ParamValues()[0]
	profile := u.controller.GetProfileByUsername(username)
	if profile == nil {
		return echo.NewHTTPError(404, "unable to find user")
	}

	currentUserId := c.Request().Context().Value("user").(uint)
	if profile.UserID != currentUserId {
		profile.FollowState = u.controller.GetFollowState(int(currentUserId), int(profile.UserID))
	}

	return c.JSON(200, profile)
}

// handle users follow and unfollow action
func (u *UsersViews) FollowView(c echo.Context) error {
	userToId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	userFromId := c.Request().Context().Value("user").(uint)
	followed := u.controller.Follow(int(userFromId), userToId)
	if !followed {
		return echo.NewHTTPError(500, "unable to follow")
	}

	return c.JSON(200, "succesfully followed")
}

func (u *UsersViews) MeView(c echo.Context) error {
	userId := c.Request().Context().Value("user").(uint)
	profile := u.controller.GetProfileById(userId)
	if profile == nil {
		return echo.NewHTTPError(404, "profile does not exist")
	}

	return c.JSON(200, profile)
}
