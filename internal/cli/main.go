package cli

import (
	"cloud-paas/internal/backend/config"
	"context"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

// Called once global flags are parsed before any subcommands are run
func rootBefore(ctx context.Context, c *cli.Command) (context.Context, error) {
	if c.Bool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Debug("Verbose output enabled")
	return ctx, nil
}

func Entrypoint() {
	config.Init()

	if err := RootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
