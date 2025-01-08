package cli

import (
	"cloud-paas/internal/backend/config"
	"context"
	"log"
	"os"
)

func Entrypoint() {
	config.Init()

	if err := RootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
