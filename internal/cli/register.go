package cli

import (
	"cloud-paas/internal/cli/config"
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func RegisterAction(ctx context.Context, c *cli.Command) error {
	conf := config.Get()
	if conf.AuthToken != "" {
		return fmt.Errorf("already logged in")
	}

	println("register called")
	return nil
}
