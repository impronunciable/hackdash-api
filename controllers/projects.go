package controllers

import (
	"app/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// InitProjectRoutes Initialize Project routes
func InitProjectRoutes(r *echo.Group) {
	r.Get("/projects", listProjects)
	r.Get("/projects/:id", getProject)
	r.Post("/projects", createProject)
	r.Patch("/projects/:id", updateProject)
	r.Delete("/projects/:id", deleteProject)
}

// Project controllers
func listProjects(c *echo.Context) error {
	projects := []models.Project{}
	models.Paginate(&models.DB, c).Find(&projects)
	return c.JSON(http.StatusOK, projects)
}

func getProject(c *echo.Context) error {
	project := models.Project{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	if models.DB.First(&project, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, project)
}

func createProject(c *echo.Context) error {
	project := models.Project{}
	dashboard := models.Dashboard{}

	if err := c.Bind(&project); err != nil {
		return err
	}

	if models.DB.First(&dashboard, project.DashboardID).RecordNotFound() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid dashboard_id parameter")
	}

	// TODO add tags
	user := c.Get("User").(models.User)
	project.UserID = user.ID
	project.Contributors = []models.User{user}
	project.Followers = []models.User{user}

	if err := models.DB.Save(&project).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, &project)
}

func updateProject(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	project := models.Project{}
	if models.DB.First(&project, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if user := c.Get("User").(models.User); user.ID != project.UserID {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	updatedData := models.Project{}
	if err := c.Bind(&updatedData); err != nil {
		return err
	}

	// TODO update tags
	record := map[string]interface{}{
		"title":       updatedData.Title,
		"description": updatedData.Description,
		"status":      updatedData.Status,
		"cover":       updatedData.Cover,
		"link":        updatedData.Link,
		"showcase":    updatedData.Showcase,
	}

	if err := models.DB.Model(&project).Updates(&record).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &project)
}

func deleteProject(c *echo.Context) error {
	project := models.Project{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id parameter missing or not a number")
	}

	if models.DB.First(&project, id).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if user := c.Get("User").(models.User); user.ID != project.UserID {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	if err := models.DB.Delete(&project).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
