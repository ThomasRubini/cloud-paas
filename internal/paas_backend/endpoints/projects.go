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
	projects := []models.DBApplication{
		{Name: "Project1", Desc: "Description1", SourceURL: "http://source1.com"},
		{Name: "Project2", Desc: "Description2", SourceURL: "http://source2.com"},
	}

	c.JSON(200, projects)
}
