package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func InitV3Routes(r *echo.Group) {
	// V3 routes definition

	// Dashboards
	r.Get("/dashboards", listDashboards)
	r.Get("/dashboards/:id", getDashboard)
	//r.Post("/dashboards", createDashboard)
	//r.Patch("/dashboards/:id", updateDashboard)
	r.Delete("/dashboards/:id", deleteDashboard)
}

// Dashboard controllers
func listDashboards(c *echo.Context) error {
	dashboards := []Dashboard{}
	Paginate(&db, c).Find(&dashboards)
	return c.JSON(http.StatusOK, dashboards)
}

func getDashboard(c *echo.Context) error {
	dashboard := Dashboard{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	if db.First(&dashboard, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, dashboard)
}

func deleteDashboard(c *echo.Context) error {
	dashboard := Dashboard{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	if db.First(&dashboard, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	db.Delete(&dashboard)
	return c.NoContent(http.StatusNoContent)
}