package clicmds

import (
	"context"
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_cli/utils"
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
	envName := cmd.Args().First()
	if envName == "" {
		return fmt.Errorf("env name is required")
	}

	env := comm.CreateEnvRequest{
		Name: envName,
	}

	resp, err := utils.GetAPIClient().R().SetPathParams(map[string]string{
		"app_id": appName,
	}).SetBody(&env).Post("/api/v1/applications/{app_id}/environments")
	if err != nil {
		return fmt.Errorf("failed to create env: %s", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to create env: %s", resp.String())
	}

	fmt.Printf("Environment %s created successfully\n", envName)
	return nil
}

func GetEnvListAction(ctx context.Context, cmd *cli.Command) error {
	var envs []comm.EnvView
	resp, err := utils.GetAPIClient().R().SetResult(&envs).SetPathParams(map[string]string{
		"app_id": appName,
	}).Get("/api/v1/applications/{app_id}/environments")
	if err != nil {
		return fmt.Errorf("failed to get env list: %s", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to get env list: %s", resp.String())
	}

	if len(envs) == 0 {
		fmt.Printf("No environments\n")
	} else {
		fmt.Printf("Environments:\n")
		for _, env := range envs {
			fmt.Printf("- %v\n", env.Name)
		}
	}

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
	envName := cmd.Args().First()
	if envName == "" {
		return fmt.Errorf("env name is required")
	}

	resp, err := utils.GetAPIClient().R().SetPathParams(map[string]string{
		"app_id": appName,
		"env_id": envName,
	}).Delete("/api/v1/applications/{app_id}/environments/{env_id}")
	if err != nil {
		return fmt.Errorf("failed to delete env: %s", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete env: %s", resp.String())
	}
	return nil
}
