package utils

import (
	"blogv2/users/models"
	"regexp"

	"github.com/labstack/echo/v4"
)

// Validate users input when they are register
func ValidateRegister(input models.UserRegisterInput) *echo.HTTPError {
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, input.Email); !m {
		return InvalidInput("email", "invalid email")
	}

	if m, _ := regexp.MatchString(`\+?1?\d{9,15}$`, input.PhoneNumber); !m {
		return InvalidInput("phoneNumber", "invalid phoneNumber")
	}

	if input.Password != input.PasswordConfirmation {
		return InvalidInput("password", "passwords do not match")
	}

	if len(input.Password) < 8 {
		return InvalidInput("password", "password must be at least 4 characts")
	}

	return nil
}

// return a bad request error with a field and
// message for the error.
func InvalidInput(field, message string) *echo.HTTPError {
	return echo.NewHTTPError(400, models.ErrorMessage{
		Field:   field,
		Message: message,
	})
}
