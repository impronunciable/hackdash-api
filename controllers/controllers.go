package controllers

import (
  "github.com/labstack/echo"
)

func InitV3Routes(r *echo.Group) {
  InitDashboardRoutes(r)
  InitUserRoutes(r)
}
