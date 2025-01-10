package cli

import (
	"cloud-paas/internal/cli/config"
	"context"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func setupLogging(c *cli.Command) {
	if c.Bool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Debug("Verbose output enabled")
}

// Called once global flags are parsed before any subcommands are run
func rootBefore(ctx context.Context, c *cli.Command) (context.Context, error) {
	setupLogging(c)

	logrus.Debugf("Server URL: %v", config.Get().BACKEND_URL)

	return ctx, nil
}

func Entrypoint() {
	config.Init()

	if err := RootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
