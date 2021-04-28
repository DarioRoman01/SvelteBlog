package main

import (
	"blogv2/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to read env: ", err)
	}

	e := echo.New()
	routes.SetRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
