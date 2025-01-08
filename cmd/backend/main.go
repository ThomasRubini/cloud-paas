package main

import (
	"fmt"
)

func main() {
	g := setupWebServer()
	launchWebServer(g)

	fmt.Println("Backend main")
}
