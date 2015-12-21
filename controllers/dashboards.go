package controllers

import (
	"app/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// InitDashboardRoutes Initialize Dashboard routes
func InitDashboardRoutes(r *echo.Group) {
	r.Get("/dashboards", listDashboards)
	r.Get("/dashboards/:id", getDashboard)
	r.Post("/dashboards", createDashboard)
	r.Patch("/dashboards/:id", updateDashboard)
	r.Delete("/dashboards/:id", deleteDashboard)
}

// Dashboard controllers
func listDashboards(c *echo.Context) error {
	dashboards := []models.Dashboard{}
	models.Paginate(&models.DB, c).Find(&dashboards)
	return c.JSON(http.StatusOK, dashboards)
}

func getDashboard(c *echo.Context) error {
	dashboard := models.Dashboard{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	if models.DB.First(&dashboard, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, dashboard)
}

func createDashboard(c *echo.Context) error {
	dashboard := models.Dashboard{}

	if err := c.Bind(&dashboard); err != nil {
		return err
	}
	dashboard.UserID = c.Get("User").(models.User).ID

	models.DB.Save(&dashboard)
	return c.JSON(http.StatusCreated, &dashboard)
}

func updateDashboard(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	dashboard := models.Dashboard{}
	if models.DB.First(&dashboard, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if user := c.Get("User").(models.User); user.ID != dashboard.UserID {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	updatedData := models.Dashboard{}
	if err := c.Bind(&updatedData); err != nil {
		return err
	}

	models.DB.Model(&dashboard).Updates(map[string]interface{}{
		"slug":        updatedData.Slug,
		"title":       updatedData.Title,
		"description": updatedData.Description,
		"link":        updatedData.Link,
		"open":        updatedData.Open,
	})
	return c.JSON(http.StatusOK, &dashboard)
}

func deleteDashboard(c *echo.Context) error {
	dashboard := models.Dashboard{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	if models.DB.First(&dashboard, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if user := c.Get("User").(models.User); user.ID != dashboard.UserID {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	models.DB.Delete(&dashboard)
	return c.NoContent(http.StatusNoContent)
}
