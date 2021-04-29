package views

import (
	"blogv2/users/controllers"
	"blogv2/users/models"
	"blogv2/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all views related with users profiles
type ProfileViews struct {
	DB *gorm.DB
}

// shortcut for not writing controllers.ProfileController
var profileController *controllers.ProfileController

// handle profile creation and validate all required fields
func (p *ProfileViews) CreateProfileView(c echo.Context) error {
	var profile models.Profile
	if err := c.Bind(&profile); err != nil {
		return echo.NewHTTPError(423, "unable to parse request body")
	}

	userId := c.Request().Context().Value("user").(uint)
	profile.UserID = userId
	httpErr := profileController.CreateProfile(&profile, p.DB)
	utils.CheckHttpError(httpErr)
	return c.JSON(201, profile)
}

// handle profile updates
func (p *ProfileViews) UpdateProfileView(c echo.Context) error {
	var profile models.Profile
	id, err := strconv.Atoi(c.Param("id"))
	utils.CheckIDParamError(err)

	err = c.Bind(&profile)
	utils.CheckRequestBodyError(err)

	userId := c.Request().Context().Value("user").(uint)
	newProfile, httpErr := profileController.UpdateProfile(userId, uint(id), &profile, p.DB)
	utils.CheckHttpError(httpErr)

	return c.JSON(200, newProfile)
}

// retrieve profile with given username
func (p *ProfileViews) GetProfileView(c echo.Context) error {
	username := c.Param("username")
	profile := profileController.GetProfileByUsername(username, p.DB)
	if profile == nil {
		return echo.NewHTTPError(404, "unable to find that user")
	}

	return c.JSON(200, profile)
}

// handle users follow and unfollow action
func (p *ProfileViews) FollowView(c echo.Context) error {
	userToId, err := strconv.Atoi(c.Param("id"))
	utils.CheckIDParamError(err)

	userFromId := c.Request().Context().Value("user").(int)
	followed := profileController.Follow(userFromId, userToId, p.DB)
	if !followed {
		return echo.NewHTTPError(500, "unable to follow")
	}

	return c.JSON(200, "succesfully followed")
}
