package endpoints

import (
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/state"

	"github.com/gin-gonic/gin"
)

func initApplications(g *gin.RouterGroup) {
	g.GET("/applications", getApps)
	g.POST("/applications", createApp)
}

// GetApplications godoc
// @Summary      List applications you have access to
// @Tags         applications
// @Produce      json
// @Success      200
// @Router       /api/v1/applications [get]
// @Success      200 {array} models.DBApplication
func getApps(c *gin.Context) {

	var apps []models.DBApplication

	if err := state.Get().Db.Find(&apps).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, apps)
}

// Postapplications godoc
// @Summary      Create a new application
// @Tags         applications
// @Accept       json
// @Param        application body models.DBApplication true "application to create"
// @Success      200
// @Router       /api/v1/applications [post]
func createApp(c *gin.Context) {
	var newApp models.DBApplication

	if err := c.ShouldBindJSON(&newApp); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := state.Get().Db.Create(&newApp).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
