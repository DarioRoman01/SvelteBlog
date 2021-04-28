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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to read env: ", err)
	}
}

var Store = sessions.NewCookieStore([]byte(os.Getenv("SECRET-KEY")))

func GenerateToken(id uint, tokenType string) (string, *echo.HTTPError) {

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
