// This package handles the webserver endpoints code
package endpoints

import (
	"github.com/gin-gonic/gin"
)

func Init(g *gin.RouterGroup) {
	initApplications(g)
	initEnvironments(g)
	initRegister(g)
}
