package backend

import (
	_ "cloud-paas/internal/backend/docs"
	"cloud-paas/internal/backend/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/swag"
)

func createWebServer() *gin.Engine {
	g := gin.New()

	g.Use(gin.Recovery())
	g.Use(middlewares.Log)

	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world !")
	})

	return g
}

func setupWebServer() *gin.Engine {
	g := createWebServer()
	setupSwag(g)
	return g
}

func setupSwag(g *gin.Engine) {
	if (swag.GetSwagger("swagger")) == nil {
		g.GET("/swagger/*any", func(c *gin.Context) {
			c.String(http.StatusNotImplemented, "OpenAPI spec was not enabled/generated")
		})
	} else {
		handler := ginSwagger.WrapHandler(swaggerFiles.Handler)
		g.GET("/swagger/*any", func(ctx *gin.Context) {
			if ctx.Request.URL.Path == "/swagger/" {
				ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			} else {
				handler(ctx)
			}
		})
	}
}

func launchWebServer(g *gin.Engine) {
	err := g.Run(":8080")
	if err != nil {
		panic(err)
	}
	panic("Web server stopped")
}
