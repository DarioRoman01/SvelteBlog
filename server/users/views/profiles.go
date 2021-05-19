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
		return utils.RequestBodyError
	}

	userId := c.Request().Context().Value("user").(uint)
	profile.UserID = userId
	httpErr := profileController.CreateProfile(&profile, p.DB)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(201, profile)
}

// handle profile updates
func (p *ProfileViews) UpdateProfileView(c echo.Context) error {
	var profile models.Profile
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	if err = c.Bind(&profile); err != nil {
		return utils.RequestBodyError
	}

	userId := c.Request().Context().Value("user").(uint)
	newProfile, httpErr := profileController.UpdateProfile(userId, uint(id), &profile, p.DB)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(200, newProfile)
}

// retrieve profile with given username
func (p *ProfileViews) GetProfileView(c echo.Context) error {
	username := c.ParamValues()[0]
	profile := profileController.GetProfileByUsername(username, p.DB)
	if profile == nil {
		return echo.NewHTTPError(404, "unable to find user")
	}

	currentUserId := c.Request().Context().Value("user").(uint)
	if profile.UserID != currentUserId {
		profile.FollowState = profileController.GetFollowState(int(currentUserId), int(profile.UserID), p.DB)
	}

	return c.JSON(200, profile)
}

// handle users follow and unfollow action
func (p *ProfileViews) FollowView(c echo.Context) error {
	userToId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.IdParamError
	}

	userFromId := c.Request().Context().Value("user").(uint)
	followed := profileController.Follow(int(userFromId), userToId, p.DB)
	if !followed {
		return echo.NewHTTPError(500, "unable to follow")
	}

	return c.JSON(200, "succesfully followed")
}

func (p *ProfileViews) MeView(c echo.Context) error {
	userId := c.Request().Context().Value("user").(uint)
	profile := profileController.GetProfileById(userId, p.DB)
	if profile == nil {
		return echo.NewHTTPError(404, "profile does not exist")
	}

	return c.JSON(200, profile)
}
