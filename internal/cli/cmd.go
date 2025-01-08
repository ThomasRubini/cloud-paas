package cli

import (
	"context"

	"github.com/urfave/cli/v3"
)

var loginCmd = &cli.Command{
	Name:  "login",
	Usage: "login with your account",
	Action: func(ctx context.Context, c *cli.Command) error {
		println("login called")
		return nil
	},
}

var registerCmd = &cli.Command{
	Name:  "register",
	Usage: "Register an account against the PaaS",
	Action: func(ctx context.Context, c *cli.Command) error {
		println("add called")
		return nil
	},
}

var projectRegisterCmd = &cli.Command{
	Name:  "register",
	Usage: "Register a new project",
	Action: func(ctx context.Context, c *cli.Command) error {
		println("project register called")
		return nil
	},
}

var projectListCmd = &cli.Command{
	Name:  "list",
	Usage: "List all projects",
	Action: func(ctx context.Context, c *cli.Command) error {
		println("project list called")
		return nil
	},
}

var projectCmd = &cli.Command{
	Commands: []*cli.Command{
		projectRegisterCmd,
		projectListCmd,
	},
}

var RootCmd = &cli.Command{
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "Enable verbose output",
			Aliases: []string{"v"},
		},
	},
	Commands: []*cli.Command{
		registerCmd,
		loginCmd,
		projectCmd,
	},

	Before: rootBefore,
}
