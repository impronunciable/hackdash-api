package main

import(
    "fmt"
    "github.com/labstack/echo"
    mw "github.com/labstack/echo/middleware"
)

func main() {
    // Create app
    app := echo.New()

    // Get config
    config := LoadConfig()

    // Basic middleware
    app.Use(mw.Logger())
    app.Use(mw.Recover())

    // Initialize router. v3 is the first version for compatibility reasons
    v3Router := app.Group("/v3")
    InitV3Routes(v3Router)

    // Initialize database
    InitDB(config.Database)

    app.Run(fmt.Sprintf(":%d", config.Port))
}
