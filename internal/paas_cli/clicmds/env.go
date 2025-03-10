package clicmds

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var appName = ""

var EnvCmd = &cli.Command{
	Name:   "env",
	Usage:  "Interact with users applications environments",
	Action: EnvCmdAction,
}

var subEnvCmd = &cli.Command{
	Name: "env",
	Commands: []*cli.Command{
		envCreateCmd,
		envListCmd,
		envInfoCmd,
		envEditCmd,
		envDeleteCmd,
	},
}

var envCreateCmd = &cli.Command{
	Name:   "create",
	Usage:  "Create an environment for given application",
	Action: createEnvAction,
}

var envListCmd = &cli.Command{
	Name:   "list",
	Usage:  "List all environments of a specific application",
	Action: GetEnvListAction,
}

var envInfoCmd = &cli.Command{
	Name:   "info",
	Usage:  "Get information about a specific environment from a given application",
	Action: GetEnvInfoAction,
}

var envEditCmd = &cli.Command{
	Name:   "edit",
	Usage:  "Edit a given environment from a given application",
	Action: editEnvAction,
}

var envDeleteCmd = &cli.Command{
	Name:   "delete",
	Usage:  "Remove a given environment from a given application",
	Action: deleteEnvAction,
}

func EnvCmdAction(ctx context.Context, cmd *cli.Command) error {
	appName = cmd.Args().First()
	return subEnvCmd.Run(ctx, cmd.Args().Slice())
}

func createEnvAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Creating env %s for application %s...\n", cmd.Args().First(), appName)
	return nil
}

func GetEnvListAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Listing all environments for application %s...\n", appName)
	return nil
}

func GetEnvInfoAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Getting env %s informations for application %s...\n", cmd.Args().First(), appName)
	return nil
}

func editEnvAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Editing env %s for application %s...\n", cmd.Args().First(), appName)
	return nil
}

func deleteEnvAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Deleting env %s for application %s...\n", cmd.Args().First(), appName)
	return nil
}
