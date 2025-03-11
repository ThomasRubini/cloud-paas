package paas_backend

import (
	"net/http"

	_ "github.com/ThomasRubini/cloud-paas/internal/paas_backend/docs"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/endpoints"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/middlewares"
	"github.com/ThomasRubini/cloud-paas/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/swag"
)

func useState(state utils.State) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("state", state)
		c.Next()
	}
}

func createWebServer(state utils.State) *gin.Engine {
	g := gin.New()

	g.Use(gin.Recovery())
	g.Use(useState(state))
	g.Use(middlewares.Log)

	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world !")
	})

	return g
}

func setupHealth(g *gin.Engine) {
	g.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}

func SetupWebServer(state utils.State) *gin.Engine {
	g := createWebServer(state)
	setupSwag(g)
	setupHealth(g)
	endpoints.Init(g.Group("/api/v1"))
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
