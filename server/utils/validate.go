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

func ValidateCursor(cursor string) *echo.HTTPError {
	cursorRegex := `[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) (2[0-3]|[01][0-9]):[0-5][0-9]`
	if m, _ := regexp.MatchString(cursorRegex, cursor); !m {
		return echo.NewHTTPError(400, "invalid cursor")
	}

	return nil
}

func CheckIDParamError(err error) *echo.HTTPError {
	if err != nil {
		return echo.NewHTTPError(400, "invalid id")
	}

	return nil
}

func CheckLimitParamError(err error) *echo.HTTPError {
	if err != nil {
		return echo.NewHTTPError(400, "invalid limit")
	}

	return nil
}

func CheckRequestBodyError(err error) *echo.HTTPError {
	if err != nil {
		return echo.NewHTTPError(423, "unable to parse request body")
	}

	return nil
}

func CheckHttpError(err *echo.HTTPError) *echo.HTTPError {
	if err != nil {
		return err
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
