package backend

import (
	"cloud-paas/internal/backend/endpoints"
	"fmt"
)

func Entrypoint() {

	g := setupWebServer()
	endpoints.Init(g.Group("/api/v1"))
	launchWebServer(g)

	fmt.Println("Backend main")
}
