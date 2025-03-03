package endpoints

import (
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"

	"github.com/gin-gonic/gin"
)

func initProjects(g *gin.RouterGroup) {
	g.GET("/projects", getProjects)
}

// GetProjects godoc
// @Summary      List projects you have access to
// @Tags         projects
// @Produce      json
// @Success      200
// @Router       /api/v1/projects [get]
// @Success      200 {array} models.DBApplication
func getProjects(c *gin.Context) {
	
	var projects []models.DBApplication
	
	if err := state.Get().Db.Find(&projects).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(200, projects)
}
