package main

import (
	"app/middleware/auth"
	"fmt"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// Create app
	app := echo.New()

	// Get config
	config := LoadConfig()

	auth0Secret, error := auth.DecodeBase64Secret(config.Auth0Secret)

	if error != nil {
		fmt.Errorf("Error decoding jwt secret: %s", error)
		panic(error)
	}

	jwtMiddleware := auth.NewJwtMiddleware("Bearer", auth0Secret, config.Auth0ClientId)

	// Basic middleware
	app.Use(mw.Logger())
	app.Use(mw.Recover())

	// Initialize router. v3 is the first version for compatibility reasons
	v3Router := app.Group("/v3")
	if !config.DevMode {
		v3Router.Use(jwtMiddleware.Handler())
	} else {
		logger.Printf("devMode is on, auth handler unregistered.")
	}
	InitV3Routes(v3Router)

	// Initialize database
	// InitDB(config.Database)

	logger.Printf("application running on port %d", config.Port)
	app.Run(fmt.Sprintf(":%d", config.Port))
}
