package endpoints

import (
	"errors"
	"strconv"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/state"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func initApplications(g *gin.RouterGroup) {
	g.GET("/applications", getApps)
	g.POST("/applications", createApp)
	g.DELETE("/applications/:app_id", deleteApp)
	g.PATCH("/applications/:app_id", updateApp)
}

// Construct an app, guessing if it's "ID" is its databsae ID or its app name
func constructAppFromId(app_id string) *models.DBApplication {
	var app models.DBApplication
	n, err := strconv.Atoi(app_id)
	if err != nil {
		app.Name = app_id
	} else {
		app.ID = uint(n)
	}
	return &app
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

	// Create in DB
	resp := state.Get().Db.Create(&newApp)
	if err := resp.Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Store credentials

	// Update credentials
	if request.SourceUsername != "" || request.SourcePassword != "" {
		err := newApp.SetSourceCredentials(request.SourceUsername, request.SourcePassword)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	logrus.Debugf("Created new application with ID %d", newApp.ID)
	c.JSON(200, IdResponse{ID: newApp.ID})
}

// DeleteApplication godoc
// @Summary      Delete an application
// @Tags         applications
// @Param        app_id path string true "Application ID"
// @Success      200
// @Router       /api/v1/applications/{app_id} [delete]
func deleteApp(c *gin.Context) {
	appConstraint := constructAppFromId(c.Param("app_id"))

	resp := state.Get().Db.Delete(&models.DBApplication{}, appConstraint)
	if err := resp.Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if resp.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "application not found"})
		return
	}

	c.Status(200)
}

// UpdateApplication godoc
// @Summary      Update an existing application
// @Tags         applications
// @Accept       json
// @Param        app_id path string true "Application ID"
// @Param        application body CreateAppRequest true "application to update"
// @Success      200
// @Router       /api/v1/applications/{app_id} [patch]
func updateApp(c *gin.Context) {
	appId := c.Param("app_id")
	if appId == "" {
		c.JSON(400, gin.H{"error": "missing id"})
		return
	}

	var app models.DBApplication
	if err := state.Get().Db.First(&app, appId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "application not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	var request CreateAppRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.CopyFields(&request, &app)

	// Update db
	db := state.Get().Db
	if err := db.Model(&app).Updates(&request).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Update credentials
	if request.SourceUsername != "" || request.SourcePassword != "" {
		err := app.SetSourceCredentials(request.SourceUsername, request.SourcePassword)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	c.Status(200)
}
