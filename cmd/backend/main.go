package main

import (
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/imgbuild"
)

func main() {
	err := imgbuild.Build("/home/itrooz/tmp2", "test")
	if err != nil {
		fmt.Println(err)
	}

	paas_backend.Entrypoint()
}
