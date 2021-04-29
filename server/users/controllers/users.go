package controllers

import (
	"blogv2/users/models"
	"blogv2/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// group all the functions related with the db and users
type UsersController struct{}

var passwordCfg *utils.PasswordConfig

func init() {
	passwordCfg = &utils.PasswordConfig{
		Time:    1,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
}

// Create user validates that the username  or email is already in use and hash users password
func (u *UsersController) CreateUser(user *models.User, db *gorm.DB) (*models.User, *echo.HTTPError) {
	emailTaken := u.GetUserByEmail(user.Email, db)
	if emailTaken != nil {
		return nil, utils.InvalidInput("email", "email already in use")
	}

	hashPassword, err := utils.GeneratePassword(passwordCfg, user.Password)
	if err != nil {
		return nil, echo.NewHTTPError(500, "unable to hash password")
	}

	user.Password = hashPassword
	if err := db.Create(&user).Error; err != nil {
		return nil, echo.NewHTTPError(500, "unable to create user")
	}

	return user, nil
}

// Login by email handles users login with email and validate that given email exist
func (u *UsersController) LoginByEmail(input *models.UserLoginInput, db *gorm.DB) (*models.User, *echo.HTTPError) {
	user := u.GetUserByEmail(input.Email, db)
	if user == nil {
		return nil, utils.InvalidInput("email", "email does not exist")
	}

	ok, _ := utils.ComparePasswords(input.Password, user.Password)
	if !ok {
		return nil, utils.InvalidInput("password", "invalid credentials")
	}

	return user, nil
}

// retrieve all users followers and paginate them
func (u *UsersController) GetUserFollowers(userId, limit int, cursor *string, db *gorm.DB) ([]models.User, bool) {
	var users []models.User

	if cursor != nil {
		db.Raw(`
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
		db.Raw(`
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
func (u *UsersController) ChangePassword(id, newPassword string, db *gorm.DB) (*models.User, *echo.HTTPError) {
	user := u.GetUserByid(id, db)
	if user == nil {
		return nil, utils.InvalidInput("token", "invalid token")
	}

	hashPwd, _ := utils.GeneratePassword(passwordCfg, newPassword)
	db.Model(&user).Update("password", hashPwd)

	return user, nil
}

// return user by email
func (u *UsersController) GetUserByEmail(email string, db *gorm.DB) *models.User {
	var user models.User
	db.Model(&models.User{}).Where("email = ?", email).Find(&user)
	if user.ID == 0 {
		return nil
	}

	return &user
}

// return user by id
func (u *UsersController) GetUserByid(id interface{}, db *gorm.DB) *models.User {
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		return nil
	}

	return &user
}
