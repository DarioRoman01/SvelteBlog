package controllers

import (
	"blogv2/users/models"
	"blogv2/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProfileController struct{}

func (p *ProfileController) CreateProfile(profile *models.Profile, db *gorm.DB) *echo.HTTPError {
	usernameTaken := p.getProfileByUsername(profile.Username, db)
	if usernameTaken {
		return utils.InvalidInput("username", "username already in use")
	}

	if err := db.Create(&profile).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create profile")
	}

	return nil
}

func (p *ProfileController) UpdateProfile(userID uint, id uint, data *models.Profile, db *gorm.DB) (*models.Profile, *echo.HTTPError) {
	var storeProfile models.Profile

	db.Model(&models.Profile{}).Where("id = ?", id).Find(&storeProfile)

	if storeProfile.ID == 0 {
		return nil, echo.NewHTTPError(404, "profile not found")
	}

	if storeProfile.UserID != userID {
		return nil, echo.NewHTTPError(403, "you do not have permissions to perform this action")
	}

	if err := db.Model(&storeProfile).Updates(data).Error; err != nil {
		return nil, echo.NewHTTPError(500, "unable to update expense")
	}

	return &storeProfile, nil
}

func (p *ProfileController) getProfileByUsername(username string, db *gorm.DB) bool {
	var profile models.Profile

	db.Model(&models.Profile{}).Where("username = ?", username).Find(&profile)

	return profile.ID != 0
}
