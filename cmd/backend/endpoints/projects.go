package endpoints

import "github.com/gin-gonic/gin"

type Project struct {
	Name       string
	Desc       string
	SourceURL  string
	AutoDeploy string
}

func initProjects(g *gin.RouterGroup) {
	g.GET("/projects", getProjects)
}

// PutModel godoc
// @Summary      List projects you have access to
// @Tags         projects
// @Produce      json
// @Success      200
// @Router       /api/v1/projects [get]
// @Success      200 {array} Project
func getProjects(c *gin.Context) {
	projects := []Project{
		{Name: "Project1", Desc: "Description1", SourceURL: "http://source1.com"},
		{Name: "Project2", Desc: "Description2", SourceURL: "http://source2.com"},
	}

	c.JSON(200, projects)
}
