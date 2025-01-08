package main

import (
	"cloud-paas/cmd/backend/endpoints"
	"fmt"
)

func main() {
	g := setupWebServer()
	endpoints.Init(g.Group("/api/v1"))
	launchWebServer(g)

	fmt.Println("Backend main")
}
