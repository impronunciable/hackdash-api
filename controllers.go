package main

import (
    "strconv"
    "net/http"
    "github.com/labstack/echo"
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
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    if db.First(&dashboard, id).RecordNotFound() {
        return echo.NewHTTPError(http.StatusNotFound)
    }
    return c.JSON(http.StatusOK, dashboard)
}

func deleteDashboard(c *echo.Context) error {
    dashboard := Dashboard{}
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    if db.First(&dashboard, id).RecordNotFound() {
        return echo.NewHTTPError(http.StatusNotFound)
    }

    db.Delete(&dashboard)

    return c.NoContent(http.StatusNoContent)
}
