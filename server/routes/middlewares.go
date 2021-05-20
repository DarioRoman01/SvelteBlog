package routes

import (
	"blogv2/utils"
	"context"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to read env: ", err)
	}
}

// routes to be skipped
var SkipperRoutes = [5]string{"/login", "/register", "/forgot-password", "/change-password", "/verify"}

/*
IsAuth middleware check if the path has to be skipped,
validate cookie and jwt token sended by the client
and set the user id in the context to get the id
of the requesting user more easyli.
*/
func IsAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if contains(SkipperRoutes, c.Path()) {
			return next(c)
		}

		session, _ := utils.Store.Get(c.Request(), "jwt")
		strToken, exists := session.Values["token"]
		if !exists {
			return echo.NewHTTPError(401, "not authenticated")
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(strToken.(string), claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT-SECRET")), nil
		})
		if err != nil {
			return echo.NewHTTPError(423, "invalid token")
		}

		if claims["type"] != utils.SESSION || !token.Valid {
			return echo.NewHTTPError(400, "invalid token")
		}

		ctx := context.WithValue(c.Request().Context(), "user", uint(claims["user_id"].(float64)))
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

// cors config middleware
func CORSconfig() echo.MiddlewareFunc {
	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("CORS_ORIGIN")},
		AllowCredentials: true,
	})

	return cors
}

func contains(arr [5]string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
