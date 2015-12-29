package main

import (
	"app/controllers"
	"app/middleware"
	"app/middleware/auth"
	"app/models"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/rs/cors"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// Create app
	app := echo.New()

	// Get config
	config := LoadConfig()

	auth0Secret, error := auth.DecodeBase64Secret(config.Auth0Secret)

	if error != nil {
		jwtErr := fmt.Errorf("Error decoding jwt secret: %s", error.Error())
		panic(jwtErr)
	}

	//Initialize database
	models.InitDB(config.Database)

	// Basic middleware
	app.Use(mw.Logger())
	app.Use(mw.Recover())
	app.Use(cors.Default().Handler)

	// Initialize router. v3 is the first version for compatibility reasons
	v3Router := app.Group("/v3")

	if !config.DevMode {

		restrict := middleware.NewRestrictMiddleware(auth0Secret, config.Auth0ClientId)
		v3Router.Use(restrict.Handler())

	} else {
		logger.Printf("devMode is on, auth handler unregistered.")
	}

	controllers.InitV3Routes(v3Router)

	logger.Printf("application running on port %d", config.Port)
	app.Run(fmt.Sprintf(":%d", config.Port))
}
