package controllers

import (
	mw "app/middleware"
	"app/models"
	"net/http"

	"github.com/labstack/echo"
)

func InitUserRoutes(r *echo.Group) {
	r.Get("/user", mw.Secure(getUser))
}

func getUser(c *echo.Context) error {
	user := c.Get("User").(models.User)
	return c.JSON(http.StatusOK, user)
}
