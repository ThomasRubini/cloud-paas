package endpoints

import "github.com/gin-gonic/gin"

func Init(g *gin.RouterGroup) {
	initProjects(g)
}
