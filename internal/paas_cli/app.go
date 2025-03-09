package paas_cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var AppCmd = &cli.Command{
	Name:  "app",
	Usage: "Interact with users applications",
	Commands: []*cli.Command{
		appCreateCmd,
		appListCmd,
		appInfoCmd,
		appDeleteCmd,
	},
}

var appCreateCmd = &cli.Command{
	Name:   "create",
	Usage:  "Create an application",
	Action: createAppAction,
}

var appListCmd = &cli.Command{
	Name:   "list",
	Usage:  "List all applications of your account",
	Action: GetAppListAction,
}

var appInfoCmd = &cli.Command{
	Name:   "info",
	Usage:  "Get information about a specific application",
	Action: GetAppInfoAction,
}

var appDeleteCmd = &cli.Command{
	Name:   "delete",
	Usage:  "Remove an application from your applications",
	Action: deleteAppAction,
}

func createAppAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Creating %s...\n", cmd.Args().First())
	return nil
}

func GetAppListAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Listing all applications...")
	return nil
}

func GetAppInfoAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Getting %s informations...\n", cmd.Args().First())
	return nil
}

func deleteAppAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Deleting %s...\n", cmd.Args().First())
	return nil
}
