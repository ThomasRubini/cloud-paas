package endpoints

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/logic"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func parseUInt(s string) uint {
	res, _ := strconv.Atoi(s)
	return uint(res)
}

func queryApp(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	appId := c.Param("app_id")
	if appId == "" {
		c.JSON(400, gin.H{"error": "missing app id"})
		c.Abort()
		return
	}

	app := models.DBApplication{}
	if err := state.Db.First(&app, constructAppFromId(appId)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "application not found"})
			c.Abort()
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
			c.Abort()
		}
	} else {
		c.Set("app", app)
		c.Next()
	}
}

func initEnvironments(g *gin.RouterGroup) {
	appRouter := g.Group("/applications/:app_id/environments")
	appRouter.Use(queryApp)
	appRouter.GET("", getEnvs)
	appRouter.POST("", createEnv)
	appRouter.DELETE("/:env_name", deleteEnv)
	appRouter.PATCH("/:env_name", updateEnv)
}

// GetEnvironments godoc
// @Summary      List environments associated to this application
// @Tags         environments
// @Produce      json
// @Success      200
// @Param		 app_id path string true "Application ID"
// @Router       /api/v1/applications/{app_id}/environments [get]
// @Success      200 {array} comm.EnvView
func getEnvs(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	app := c.MustGet("app").(models.DBApplication)

	var envs []models.DBEnvironment

	if err := state.Db.Where("application_id = ?", app.ID).Find(&envs).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	envViews := make([]comm.EnvView, 0)
	for _, env := range envs {
		var envView comm.EnvView
		utils.CopyFields(&env, &envView)
		envViews = append(envViews, envView)
	}

	c.JSON(200, envViews)
}

// CreateEnvironment godoc
// @Summary      Create a new environment
// @Tags         environments
// @Accept       json
// @Param		 app_id path string true "Application ID"
// @Param        environment body comm.CreateEnvRequest true "environment to create"
// @Success      200 {object} comm.IdResponse
// @Router       /api/v1/applications/{app_id}/environments/ [post]
func createEnv(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	app := c.MustGet("app").(models.DBApplication)

	var request comm.CreateEnvRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verify app does not already exist
	var count int64
	if err := state.Db.Model(&models.DBEnvironment{}).Where("application_id = ? and name = ?", app.ID, request.Name).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(409, gin.H{"error": "application already exists"})
		return
	}

	var newEnv models.DBEnvironment
	newEnv.ApplicationID = app.ID
	utils.CopyFields(&request, &newEnv)

	// Create in DB
	resp := state.Db.Create(&newEnv)
	if err := resp.Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Create deployment
	err := logic.HandleEnvironmentUpdate(state, app, newEnv)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Errorf("failed to deploy environment: %w", err).Error()})
		return
	}

	logrus.Debugf("Created new environment with ID %d", newEnv.ID)
	c.JSON(200, comm.IdResponse{ID: newEnv.ID})
}

func getDBEnv(c *gin.Context) (*models.DBEnvironment, error) {
	state := c.MustGet("state").(utils.State)
	app := c.MustGet("app").(models.DBApplication)

	envName := c.Param("env_name")
	if envName == "" {
		return nil, nil
	}

	var env models.DBEnvironment
	if err := state.Db.Where("application_id = ? AND name = ?", app.ID, envName).First(&env).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &env, nil
}

// DeleteEnvironment godoc
// @Summary      Delete an environment
// @Tags         environments
// @Param		 app_id path string true "Application ID"
// @Param        env_name path string true "Environment name"
// @Success      200
// @Router       /api/v1/applications/{app_id}/environments/{env_name} [delete]
func deleteEnv(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	env, err := getDBEnv(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else if env == nil {
		c.JSON(404, gin.H{"error": "environment does not exist"})
		return
	}

	envToDelete := models.DBEnvironment{
		ApplicationID: parseUInt(c.Param("app_id")),
		Name:          c.Param("env_name"),
	}

	if err := state.Db.Delete(&models.DBEnvironment{}, &envToDelete).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

// UpdateEnvironment godoc
// @Summary      Update an existing environment
// @Tags         environments
// @Accept       json
// @Param		 app_id path string true "Application ID"
// @Param        env_name path string true "Environment name"
// @Param        environment body comm.CreateEnvRequest true "environment to update"
// @Success      200
// @Router       /api/v1/applications/{app_id}/environments/{env_name} [patch]
func updateEnv(c *gin.Context) {
	state := c.MustGet("state").(utils.State)
	app := c.MustGet("app").(models.DBApplication)
	env, err := getDBEnv(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else if env == nil {
		c.JSON(404, gin.H{"error": "environment not found"})
		return
	}

	var request comm.CreateEnvRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.CopyFields(&request, env)

	// Update db
	if err := state.Db.Model(&env).Updates(&request).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Update deployment
	err = logic.HandleEnvironmentUpdate(state, app, *env)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Errorf("failed to deploy environment: %w", err).Error()})
		return
	}

	c.Status(200)
}
