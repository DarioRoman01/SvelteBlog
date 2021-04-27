package views

import (
	"blogv2/users/controllers"
	"blogv2/users/models"
	"blogv2/utils"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProfileViews struct {
	DB *gorm.DB
}

var profileController *controllers.ProfileController

func (p *ProfileViews) CreateProfileView(c echo.Context) error {
	var profile models.Profile

	if err := c.Bind(&profile); err != nil {
		return echo.NewHTTPError(423, "unable to parse request body")
	}

	userId, _ := strconv.Atoi(utils.UserIDFromToken(c))
	profile.UserID = uint(userId)

	httpErr := profileController.CreateProfile(&profile, p.DB)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(201, profile)
}

func (p *ProfileViews) UpdateProfileView(c echo.Context) error {
	var profile models.Profile
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.Bind(&profile); err != nil {
		return echo.NewHTTPError(423, "unable to parse request body")
	}

	userId, _ := strconv.Atoi(utils.UserIDFromToken(c))
	newProfile, httpErr := profileController.UpdateProfile(uint(userId), uint(id), &profile, p.DB)
	if httpErr != nil {
		return c.JSON(httpErr.Code, httpErr.Message)
	}

	return c.JSON(200, newProfile)
}

func (p *ProfileViews) GetProfileView(c echo.Context) error {
	username := c.Param("username")
	profile := profileController.GetProfileByUsername(username, p.DB)
	if profile == nil {
		return echo.NewHTTPError(404, "unable to find that user")
	}

	return c.JSON(200, profile)
}
