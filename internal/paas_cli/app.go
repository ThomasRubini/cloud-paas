package paas_cli

import (
	"github.com/urfave/cli/v3"
)

var AppCmd = &cli.Command{
	Name:   "app",
	Usage:  "Interact with users applications",
	Action: LoginAction,
	Commands: []*cli.Command{
		createAppCmd,
		listAppCmd,
		infoAppCmd,
		deleteAppCmd,
	},
}

var createAppCmd = &cli.Command{
	Name:   "create",
	Usage:  "Create an application",
	Action: LoginAction,
}

var listAppCmd = &cli.Command{
	Name:      "list",
	Usage:     "List all applications of your account",
	Action:    LoginAction,
	ArgsUsage: "[app]",
}

var infoAppCmd = &cli.Command{
	Name:      "info",
	Usage:     "Get information about a specific application",
	Action:    LoginAction,
	ArgsUsage: "[app]",
}

var deleteAppCmd = &cli.Command{
	Name:      "info",
	Usage:     "Remove an application from your applications",
	Action:    LoginAction,
	ArgsUsage: "[app]",
}
