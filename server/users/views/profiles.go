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
