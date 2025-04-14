package clicmds

import (
	"context"
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_cli/utils"
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
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "desc",
			Usage: "Description of the application",
		},
		&cli.StringFlag{
			Name:     "source-url",
			Usage:    "Source URL of the application",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "source-username",
			Usage: "Username of user on the repository source",
		},
		&cli.StringFlag{
			Name:  "source-password",
			Usage: "Password of user on the repository source. On Github, this can be a personal access token",
		},
	},
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
	appName := cmd.Args().First()
	if appName == "" {
		return fmt.Errorf("app name is required")
	}

	app := comm.CreateAppRequest{
		Name:           appName,
		Desc:           cmd.String("desc"),
		SourceURL:      cmd.String("source-url"),
		SourceUsername: cmd.String("source-username"),
		SourcePassword: cmd.String("source-password"),
	}

	var parsedResp comm.IdResponse
	resp, err := utils.GetAPIClient().R().SetBody(&app).SetResult(&parsedResp).Post("/api/v1/applications")
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to create app: %s", resp.String())
	}

	fmt.Printf("Application %s created successfully (id: %v)\n", appName, parsedResp.ID)
	return nil
}

func GetAppListAction(ctx context.Context, cmd *cli.Command) error {
	var apps []comm.AppView
	resp, err := utils.GetAPIClient().R().SetResult(&apps).Get("/api/v1/applications")
	if err != nil {
		return fmt.Errorf("failed to get app list: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to get app list: %s", resp.String())
	}

	if len(apps) == 0 {
		fmt.Printf("No applications\n")
	} else {
		fmt.Printf("Applications:\n")
		for _, app := range apps {
			fmt.Printf("- %v (ID: %v)\n", app.Name, app.ID)
		}
	}
	return nil
}

func GetAppInfoAction(ctx context.Context, cmd *cli.Command) error {
	var app comm.AppView
	resp, err := utils.GetAPIClient().R().SetPathParams(map[string]string{
		"app_id": cmd.Args().First(),
	}).SetResult(&app).Get("/api/v1/applications/{app_id}")
	if err != nil {
		return fmt.Errorf("failed to get app info: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to get app info: %s", resp.String())
	}

	fmt.Printf("Application info:\n")
	fmt.Printf("ID: %v\n", app.ID)
	fmt.Printf("Name: %v\n", app.Name)
	fmt.Printf("Description: %v\n", app.Desc)
	fmt.Printf("Source URL: %v\n", app.SourceURL)
	fmt.Printf("Auto Deploy: %v\n", app.AutoDeploy)

	return nil
}

func deleteAppAction(ctx context.Context, cmd *cli.Command) error {
	resp, err := utils.GetAPIClient().R().SetPathParams(map[string]string{
		"app_id": cmd.Args().First(),
	}).Delete("/api/v1/applications/{app_id}")
	if err != nil {
		return fmt.Errorf("failed to delete app: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete app: %s", resp.String())
	}

	fmt.Printf("Application deleted successfully\n")
	return nil
}
