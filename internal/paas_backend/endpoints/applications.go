package endpoints

import (
	"errors"
	"strconv"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func initApplications(g *gin.RouterGroup) {
	g.GET("/applications", getApps)
	g.GET("/applications/:app_id", getApp)
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

// GetApplications godoc
// @Summary      List applications you have access to
// @Tags         applications
// @Produce      json
// @Success      200
// @Router       /api/v1/applications [get]
// @Success      200 {array} comm.AppView
func getApps(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	var apps []models.DBApplication

	if err := state.Db.Find(&apps).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	appsViews := make([]comm.AppView, 0)
	for _, app := range apps {
		var appView comm.AppView
		utils.CopyFields(&app, &appView)
		appsViews = append(appsViews, appView)
	}

	c.JSON(200, appsViews)
}

func getApp(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	appConstraint := constructAppFromId(c.Param("app_id"))

	var app models.DBApplication
	if err := state.Db.First(&app, appConstraint).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "application not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	var appView comm.AppView
	utils.CopyFields(&app, &appView)

	c.JSON(200, appView)
}

// Postapplications godoc
// @Summary      Create a new application
// @Tags         applications
// @Accept       json
// @Param        application body comm.CreateAppRequest true "application to create"
// @Success      200 {object} comm.IdResponse
// @Router       /api/v1/applications [post]
func createApp(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	var request comm.CreateAppRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verify app does not already exist
	var count int64
	if err := state.Db.Model(&models.DBApplication{}).Where("name = ?", request.Name).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(409, gin.H{"error": "application already exists"})
		return
	}

	var newApp models.DBApplication
	utils.CopyFields(&request, &newApp)

	// Create in DB
	resp := state.Db.Create(&newApp)
	if err := resp.Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Store credentials
	if request.SourceUsername != "" || request.SourcePassword != "" {
		err := state.SecretsProvider.SetSourceCredentials(newApp, request.SourceUsername, request.SourcePassword)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	logrus.Debugf("Created new application with ID %d", newApp.ID)
	c.JSON(200, comm.IdResponse{ID: newApp.ID})
}

// DeleteApplication godoc
// @Summary      Delete an application
// @Tags         applications
// @Param        app_id path string true "Application ID"
// @Success      200
// @Router       /api/v1/applications/{app_id} [delete]
func deleteApp(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	appConstraint := constructAppFromId(c.Param("app_id"))

	resp := state.Db.Delete(&models.DBApplication{}, appConstraint)
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
// @Param        application body comm.CreateAppRequest true "application to update"
// @Success      200
// @Router       /api/v1/applications/{app_id} [patch]
func updateApp(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	appId := c.Param("app_id")
	if appId == "" {
		c.JSON(400, gin.H{"error": "missing id"})
		return
	}

	var app models.DBApplication
	if err := state.Db.First(&app, appId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "application not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	var request comm.CreateAppRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.CopyFields(&request, &app)

	// Update db
	if err := state.Db.Model(&app).Updates(&request).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Update credentials
	if request.SourceUsername != "" || request.SourcePassword != "" {
		err := state.SecretsProvider.SetSourceCredentials(app, request.SourceUsername, request.SourcePassword)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	c.Status(200)
}
