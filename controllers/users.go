package controllers

import (
  "app/models"
  "github.com/labstack/echo"
  "net/http"
)

func InitUserRoutes(r *echo.Group) {
  r.Get("/user", getUser)
}

func getUser(c *echo.Context) error {
    user := c.Get("User").(models.User)
    return c.JSON(http.StatusOK, user)
}