package main

import (
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/imgbuild"
)

func main() {
	err := imgbuild.Build("/home/itrooz/tmp2", []string{"test"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Hey")
	// paas_backend.Entrypoint()
}
