package controllers

import (
	"blogv2/users/models"
	"blogv2/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProfileController struct{}

// handle profile creation
func (p *ProfileController) CreateProfile(profile *models.Profile, db *gorm.DB) *echo.HTTPError {
	usernameTaken := p.GetProfileByUsername(profile.Username, db)
	if usernameTaken != nil {
		return utils.InvalidInput("username", "username already in use")
	}

	if err := db.Create(&profile).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create profile")
	}

	return nil
}

// handle profile update and validate that the requesting user is owner of the profile.
func (p *ProfileController) UpdateProfile(userID uint, id uint, data *models.Profile, db *gorm.DB) (*models.Profile, *echo.HTTPError) {
	var storeProfile models.Profile

	db.Model(&models.Profile{}).Where("id = ?", id).Find(&storeProfile)

	if storeProfile.UserID == 0 {
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

// retrieve profile by given username
func (p *ProfileController) GetProfileByUsername(username string, db *gorm.DB) *models.Profile {
	var profile models.Profile

	db.Model(&models.Profile{}).Where("username = ?", username).Find(&profile)

	if profile.UserID == 0 {
		return nil
	}

	return &profile
}

// handle users folow and unfollow checking if the requesting user
// is already following the given user
func (p *ProfileController) Follow(userFromID, userToID int, db *gorm.DB) bool {
	var userFrom models.User
	db.Table("users").Where("id = ?", userFromID).Find(&userFrom)
	if userFrom.ID == 0 {
		return false
	}

	var userTo models.User
	db.Table("users").Where("id = ?", userFromID).Find(&userTo)
	if userTo.ID == 0 {
		return false
	}

	following := p.GetFollowState(userFromID, userToID, db)

	if !following {
		err := db.Model(&userFrom).Association("Follow").Append(&userTo)
		if err != nil {
			return false
		}

		db.Model(&userTo.Profile).Update("followers = followers + ?", 1)
		return true
	}

	err := db.Model(&userFrom).Association("Follow").Delete(&userFrom)
	if err != nil {
		return false
	}

	db.Model(&userTo.Profile).Update("followers = followers - ?", 1)
	return true
}

// get follow state of the requesting user
func (p *ProfileController) GetFollowState(userFromID, userToID int, db *gorm.DB) bool {
	var followId int
	db.Raw(`
		SELECT follow_id,
		FROM "user_follow",
		WHERE user_id = ?
		AND follow_id =?
	`, userFromID, userToID).Find(&followId)

	return followId != 0
}
