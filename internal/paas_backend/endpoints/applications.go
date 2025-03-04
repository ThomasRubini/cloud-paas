package endpoints

import (
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/state"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func initApplications(g *gin.RouterGroup) {
	g.GET("/applications", getApps)
	g.POST("/applications", createApp)
	g.DELETE("/applications/:id", deleteApp)
}

type AppView struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	SourceURL  string `json:"source_url"`
	AutoDeploy bool   `json:"auto_deploy"`
}

// GetApplications godoc
// @Summary      List applications you have access to
// @Tags         applications
// @Produce      json
// @Success      200
// @Router       /api/v1/applications [get]
// @Success      200 {array} AppView
func getApps(c *gin.Context) {

	var apps []models.DBApplication

	if err := state.Get().Db.Find(&apps).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var appsViews []AppView
	for _, app := range apps {
		var appView AppView
		utils.CopyFields(&app, &appView)
		appsViews = append(appsViews, appView)
	}

	c.JSON(200, appsViews)
}

type IdResponse struct {
	ID uint `json:"id"`
}

type CreateAppRequest struct {
	Name           string
	Desc           string
	SourceURL      string
	SourceUsername string
	SourcePassword string
	AutoDeploy     bool
}

// Postapplications godoc
// @Summary      Create a new application
// @Tags         applications
// @Accept       json
// @Param        application body CreateAppRequest true "application to create"
// @Success      200 {object} IdResponse
// @Router       /api/v1/applications [post]
func createApp(c *gin.Context) {
	var request CreateAppRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var newApp models.DBApplication
	utils.CopyFields(&request, &newApp)

	resp := state.Get().Db.Create(&newApp)
	if err := resp.Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	logrus.Debugf("Created new application with ID %d", newApp.ID)
	c.JSON(200, IdResponse{ID: newApp.ID})
}

// DeleteApplication godoc
// @Summary      Delete an application
// @Tags         applications
// @Param        id path int true "Application ID"
// @Success      200
// @Router       /api/v1/applications/{id} [delete]
func deleteApp(c *gin.Context) {
	appId := c.Param("id")
	if appId == "" {
		c.JSON(400, gin.H{"error": "missing id"})
		return
	}

	var app models.DBApplication
	if err := state.Get().Db.First(&app, appId).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := state.Get().Db.Delete(&app).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
