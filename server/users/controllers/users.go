package controllers

import (
	"blogv2/users/models"
	"blogv2/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all the functions related with the db and users
type UsersController struct {
	db  *gorm.DB
	cfg *utils.PasswordConfig
}

func NewUserController(db *gorm.DB) *UsersController {
	return &UsersController{
		db: db,
		cfg: &utils.PasswordConfig{
			Time:    1,
			Memory:  64 * 1024,
			Threads: 4,
			KeyLen:  32,
		},
	}
}

// Create user validates that the username  or email is already in use and hash users password
func (u *UsersController) CreateUser(user *models.User) *echo.HTTPError {
	emailTaken := u.GetUserByEmail(user.Email)
	if emailTaken != nil {
		return utils.InvalidInput("email", "email already in use")
	}

	hashPassword, err := utils.GeneratePassword(u.cfg, user.Password)
	if err != nil {
		return echo.NewHTTPError(500, "unable to hash password")
	}

	user.Password = hashPassword
	if err := u.db.Table("users").Create(user).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create user")
	}

	return nil
}

// Login by email handles users login with email and validate that given email exist
func (u *UsersController) LoginByEmail(input *models.UserLoginInput) (*models.User, *echo.HTTPError) {
	user := u.GetUserByEmail(input.Email)
	if user == nil {
		return nil, utils.InvalidInput("email", "email does not exist")
	} else if !user.IsVerified {
		return nil, utils.InvalidInput("user", "your account is not active yet")
	}

	ok, _ := utils.ComparePasswords(input.Password, user.Password)
	if !ok {
		return nil, utils.InvalidInput("password", "invalid credentials")
	}

	return user, nil
}

// retrieve all users followers and paginate them
func (u *UsersController) GetUserFollowers(userId, limit int, cursor string) ([]models.User, bool) {
	var users []models.User

	if cursor != "" {
		u.db.Raw(`
			SELECT u.*,
			FROM "users"
			WHERE id in (
				SELECT follow_id 
				FROM "user_follow"
				WHERE user_id = ?
			)
			AND u.created_at < ?
			ORDER BY u.created_at DESC
			LIMIT ?
		`, userId, cursor, limit).Find(&users)
	} else {
		u.db.Raw(`
			SELECT u.*,
			FROM "users"
			WHERE id in (
				SELECT follow_id 
				FROM "user_follow"
				WHERE user_id = ?
			)
			ORDER BY u.created_at DESC
			LIMIT ?
		`, userId, limit).Find(&users)
	}

	if len(users) == 0 {
		return nil, false
	}
	if len(users) == limit {
		return users[0 : limit-1], true
	}

	return users[0 : len(users)-1], false
}

// Change Password change user password
func (u *UsersController) ChangePassword(id, newPassword string) (*models.User, *echo.HTTPError) {
	user := u.GetUserByid(id)
	if user == nil {
		return nil, utils.InvalidInput("token", "invalid token")
	}

	hashPwd, _ := utils.GeneratePassword(u.cfg, newPassword)
	u.db.Model(&user).Update("password", hashPwd)

	return user, nil
}

// return user by email
func (u *UsersController) GetUserByEmail(email string) *models.User {
	var user models.User
	u.db.Table("users").Where("email = ?", email).Find(&user)
	if user.ID == 0 {
		return nil
	}

	return &user
}

// return user by id
func (u *UsersController) GetUserByid(id interface{}) *models.User {
	var user models.User
	u.db.First(&user, id)
	if user.ID == 0 {
		return nil
	}

	return &user
}

func (u *UsersController) CreateProfile(profile *models.Profile) *echo.HTTPError {
	usernameTaken := u.GetProfileByUsername(profile.Username)
	if usernameTaken != nil {
		return utils.InvalidInput("username", "username already in use")
	}

	if err := u.db.Table("profiles").Create(&profile).Error; err != nil {
		return echo.NewHTTPError(500, "unable to create profile")
	}

	return nil
}

// handle profile update and validate that the requesting user is owner of the profile.
func (u *UsersController) UpdateProfile(userID uint, id uint, data *models.Profile) (*models.Profile, *echo.HTTPError) {
	var storeProfile models.Profile
	u.db.Table("profiles").Where("user_id = ?", id).Find(&storeProfile)

	if storeProfile.UserID == 0 {
		return nil, echo.NewHTTPError(404, "profile not found")
	}

	if storeProfile.UserID != userID {
		return nil, echo.NewHTTPError(403, "you do not have permissions to perform this action")
	}

	// only check for taken username username
	// if username is updated
	if storeProfile.Username != data.Username {
		usernameTaken := u.GetProfileByUsername(data.Username)
		if usernameTaken != nil {
			return nil, echo.NewHTTPError(400, "username already taken")
		}
	}

	if err := u.db.Model(&storeProfile).Updates(data).Error; err != nil {
		return nil, echo.NewHTTPError(500, "unable to update expense")
	}

	return &storeProfile, nil
}

// retrieve profile by given username
func (u *UsersController) GetProfileByUsername(username string) *models.Profile {
	var profile models.Profile
	u.db.Table("profiles").Where("username = ?", username).Find(&profile)
	if profile.UserID == 0 {
		return nil
	}

	return &profile
}

func (u *UsersController) GetProfileById(id uint) *models.Profile {
	var profile models.Profile
	u.db.Table("profiles").Where("user_id = ?", id).Find(&profile)
	if profile.UserID == 0 {
		return nil
	}

	return &profile
}

// handle users folow and unfollow checking if the requesting user
// is already following the given user
func (u *UsersController) Follow(userFromID, userToID int) bool {
	var userFrom models.User
	u.db.Table("users").Where("id = ?", userFromID).Find(&userFrom)
	if userFrom.ID == 0 {
		return false
	}

	var userTo models.User
	u.db.Table("users").Where("id = ?", userToID).Find(&userTo)
	if userTo.ID == 0 {
		return false
	}

	following := u.GetFollowState(userFromID, userToID)

	if !following {
		err := u.db.Model(&userFrom).Association("Follow").Append(&userTo)
		if err != nil {
			return false
		}

		u.db.Exec(`UPDATE "profiles" SET followers = followers + 1 WHERE user_id = ?`, userToID)
		return true
	} else {
		err := u.db.Model(&userFrom).Association("Follow").Delete(&userTo)
		if err != nil {
			return false
		}

		u.db.Exec(`UPDATE "profiles" SET followers = followers - 1 WHERE user_id = ?`, userToID)
		return true
	}
}

// get follow state of the requesting user
func (u *UsersController) GetFollowState(userFromID, userToID int) bool {
	var followId int
	u.db.Raw(`
		SELECT follow_id
		FROM "user_follow"
		WHERE user_id = ?
		AND follow_id =?
	`, userFromID, userToID).Find(&followId)

	return followId != 0
}

func (u *UsersController) UpdateUserStatus(id int) *echo.HTTPError {
	err := u.db.Table("users").Where("id = ?", id).Update("is_verified", true).Error
	if err != nil {
		return echo.NewHTTPError(500, "unable to update user status")
	}

	return nil
}
