package controllers

import (
	"app/models"
	"net/http"

	"github.com/labstack/echo"
)

// InitDashboardRoutes Initialize Dashboard routes
func InitDashboardRoutes(r *echo.Group) {
	r.Get("/dashboards", listDashboards)
	r.Get("/dashboards/:slug", getDashboard)
	r.Post("/dashboards", createDashboard)
	r.Patch("/dashboards/:slug", updateDashboard)
	r.Delete("/dashboards/:slug", deleteDashboard)
}

// Dashboard controllers
func listDashboards(c *echo.Context) error {
	dashboards := []models.Dashboard{}
	models.Paginate(&models.DB, c).Find(&dashboards)
	return c.JSON(http.StatusOK, dashboards)
}

func getDashboard(c *echo.Context) error {
	dashboard := models.Dashboard{}
	slug := c.Param("slug")

	if slug == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "slug parameter missing")
	}

	if models.DB.Preload("Projects").Where("slug = ?", slug).First(&dashboard).RecordNotFound() {
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

	if err := models.DB.Save(&dashboard).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, &dashboard)
}

func updateDashboard(c *echo.Context) error {
	slug := c.Param("slug")

	if slug == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "slug parameter missing")
	}

	dashboard := models.Dashboard{}
	if models.DB.Where("slug = ?", slug).First(&dashboard).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if user := c.Get("User").(models.User); user.ID != dashboard.UserID {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	updatedData := models.Dashboard{}
	if err := c.Bind(&updatedData); err != nil {
		return err
	}

	record := map[string]interface{}{
		"title":       updatedData.Title,
		"description": updatedData.Description,
		"link":        updatedData.Link,
		"open":        updatedData.Open,
	}

	if err := models.DB.Model(&dashboard).Updates(&record).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &dashboard)
}

func deleteDashboard(c *echo.Context) error {
	dashboard := models.Dashboard{}
	slug := c.Param("slug")

	if slug == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "slug parameter missing")
	}

	if models.DB.Where("slug = ?", slug).First(&dashboard).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if user := c.Get("User").(models.User); user.ID != dashboard.UserID {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	if err := models.DB.Delete(&dashboard).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
