package main

import (
	"blogv2/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	routes.SetRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
