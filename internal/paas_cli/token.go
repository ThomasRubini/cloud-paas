package paas_cli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var TokenCmd = &cli.Command{
	Name:  "token",
	Usage: "Interact with users tokens",
	Commands: []*cli.Command{
		tokenGenerateCmd,
		tokenListCmd,
		tokenDeleteCmd,
	},
}

var tokenGenerateCmd = &cli.Command{
	Name:   "generate",
	Usage:  "Generate a new token",
	Action: generateTokenAction,
}

var tokenListCmd = &cli.Command{
	Name:   "list",
	Usage:  "List all tokens of your account",
	Action: listTokenAction,
}

var tokenDeleteCmd = &cli.Command{
	Name:   "delete",
	Usage:  "Remove a token from your account",
	Action: deleteTokenAction,
}

func generateTokenAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Generating token...")
	return nil
}

func listTokenAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Listing all tokens...")
	return nil
}

func deleteTokenAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Deleting token %s...\n", cmd.Args().First())
	return nil
}
