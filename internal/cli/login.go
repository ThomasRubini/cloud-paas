package cli

import (
	"cloud-paas/internal/cli/config"
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func LoginAction(ctx context.Context, c *cli.Command) error {
	conf := config.Get()
	if conf.AUTH_TOKEN != "" {
		return fmt.Errorf("already logged in")
	}

	println("login called")
	return nil
}
