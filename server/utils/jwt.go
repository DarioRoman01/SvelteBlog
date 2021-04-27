package utils

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to read env: ", err)
	}
}

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

func GetToken(c echo.Context) (*jwt.Token, jwt.MapClaims) {
	headerToken := c.Request().Header.Get("x-auth-token")
	strToken := strings.Split(headerToken, " ")[1]
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(strToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT-SECRET")), nil
	})
	if err != nil {
		return nil, nil
	}

	return token, claims
}

func UserIDFromToken(c echo.Context) string {
	_, claims := GetToken(c)
	return claims["user_id"].(string)
}
