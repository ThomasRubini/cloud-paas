package clicmds

import (
	"context"
	"fmt"
	"os"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_cli/utils"
	"github.com/urfave/cli/v3"
)

var appName = ""

var EnvCmd = &cli.Command{
	Name:            "env",
	Usage:           "Interact with users applications environments",
	Action:          EnvCmdAction,
	SkipFlagParsing: true,
}

var subEnvCmd = &cli.Command{
	Name: "<app_name>",
	Commands: []*cli.Command{
		envCreateCmd,
		envListCmd,
		envInfoCmd,
		envEditCmd,
		envDeleteCmd,
		envVarsCmd,
	},
}

var envCreateCmd = &cli.Command{
	Name:   "create",
	Usage:  "Create an environment for given application",
	Action: createEnvAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "branch",
			Usage:    "Branch to use for the environment",
			Required: true,
			Aliases:  []string{"b"},
		},
		&cli.StringFlag{
			Name:     "domain",
			Usage:    "Domain to use for the environment",
			Required: true,
			Aliases:  []string{"d"},
		},
	},
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
	Usage:  "Edit environment variable from a given environment of a given application",
	Action: editEnvAction,
}

var envVarsCmd = &cli.Command{
	Name:   "vars",
	Usage:  "Edit environment variables of a given environment of a given application",
	Action: editEnvVarsAction,
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
		Name:   envName,
		Branch: cmd.String("branch"),
		Domain: cmd.String("domain"),
	}

	resp, err := utils.GetAPIClient().R().SetPathParams(map[string]string{
		"app_id": appName,
	}).SetBody(&env).Post("/api/v1/applications/{app_id}/environments")
	if err != nil {
		return fmt.Errorf("failed to create env: %w", err)
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
		return fmt.Errorf("failed to get env list: %w", err)
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

// Action for the "cli env <app_name> edit <env_name>" command
func editEnvVarsAction(ctx context.Context, cmd *cli.Command) error {
	envName := cmd.Args().First()
	if envName == "" {
		return fmt.Errorf("env name is required")
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("ENV_VARS_%s_*.yaml", envName))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	resp, err := utils.GetAPIClient().R().SetPathParams(map[string]string{
		"app_name": appName,
		"env_name": envName,
	}).SetResult(&comm.EnvView{}).Get("/api/v1/applications/{app_name}/environments/{env_name}")

	if err != nil {
		return fmt.Errorf("failed to get env: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to get env: %s", resp.String())
	}
	yamlBytes, err := utils.JSONtoYAML([]byte(resp.Result().(*comm.EnvView).EnvVars))
	if err != nil {
		return fmt.Errorf("failed to convert JSON to YAML: %w", err)
	}
	envVars := string(yamlBytes)
	envVars = fmt.Sprintf("# Add environment variables for your environment %s here in YAML format\n%s", envName, envVars)

	// Write the environment variables to the temp file
	if _, err := tempFile.WriteString(envVars); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tempFile.Close()

	// Open the temp file in the default editor
	updatedEnvVars, err := utils.OpenInEditor(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	updatedJson, err := utils.YAMLtoJSON(updatedEnvVars)
	if err != nil {
		return fmt.Errorf("failed to parse user input: %w", err)
	}
	if string(updatedJson) == resp.Result().(*comm.EnvView).EnvVars {
		fmt.Print("No changes made to environment variables.")
	} else {
		// API Call to update the environment variables
		resp, err = utils.GetAPIClient().R().SetPathParams(map[string]string{
			"app_id": appName,
			"env_id": envName,
		}).SetBody(&comm.EnvView{
			EnvVars: string(updatedJson),
		}).Patch("/api/v1/applications/{app_id}/environments/{env_id}")
		if err != nil {
			return fmt.Errorf("failed to update env: %w", err)
		}
		if resp.StatusCode() != 200 {
			return fmt.Errorf("failed to update env: %s", resp.String())
		}
		fmt.Printf("Environment %s updated successfully\n", envName)
	}
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
		return fmt.Errorf("failed to delete env: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete env: %s", resp.String())
	}

	fmt.Println("Environment deleted successfully")
	return nil
}
