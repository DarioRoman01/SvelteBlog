package utils

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET-KEY")))

type TokenType int

const (
	head = iota
	CHANGE
	VERIFY
	SESSION
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to read env: ", err)
	}
}

// Set the jwt token in the current request
func SetToken(c echo.Context, token string) {
	session, _ := store.Get(c.Request(), "jwt")
	session.Values["token"] = token
	session.Save(c.Request(), c.Response().Writer)
}

// Set the jwt token as expired
func Unabletoken(c echo.Context) error {
	session, _ := store.Get(c.Request(), "jwt")
	session.Options.MaxAge = -1
	err := session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return echo.NewHTTPError(500, "unable to delete session")
	}

	return nil
}

// check that the tokens is in the request
func CheckToken(c echo.Context) (string, error) {
	session, _ := store.Get(c.Request(), "jwt")
	strToken, exists := session.Values["token"]
	if !exists {
		return "", echo.NewHTTPError(401, "not authenticated")
	}

	return strToken.(string), nil
}

// generates new jwt token with the user id and token type
// there is only tree types "change-password",
// "verify", and "session"
func GenerateToken(id uint, tokenType TokenType) (string, *echo.HTTPError) {

	claims := jwt.MapClaims{}
	claims["user_id"] = id
	claims["type"] = tokenType
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString([]byte(os.Getenv("JWT-SECRET")))
	if err != nil {
		return "", echo.NewHTTPError(500, "Unable to create token")
	}

	return token, nil
}
