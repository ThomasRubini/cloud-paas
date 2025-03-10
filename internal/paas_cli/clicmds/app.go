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
	app := comm.CreateAppRequest{
		Name: cmd.Args().First(),
		Desc: cmd.String("desc"),
	}

	var parsedResp comm.IdResponse
	resp, err := utils.GetAPIClient().R().SetBody(&app).SetResult(&parsedResp).Post("/api/v1/applications")
	if err != nil {
		fmt.Printf("Failed to create app: %s\n", err)
	}

	if resp.StatusCode() != 200 {
		fmt.Printf("Failed to create app: %s\n", resp.String())
	}

	fmt.Printf("Application created successfully (id: %v)\n", parsedResp.ID)
	return nil
}

func GetAppListAction(ctx context.Context, cmd *cli.Command) error {
	var apps []comm.AppView
	resp, err := utils.GetAPIClient().R().SetResult(&apps).Get("/api/v1/applications")
	if err != nil {
		fmt.Printf("Failed to get app list: %s\n", err)
	}

	if resp.StatusCode() != 200 {
		fmt.Printf("Failed to get app list: %s\n", resp.String())
	}

	fmt.Printf("Applications:\n")
	for _, app := range apps {
		fmt.Printf("- %v (ID: %v)\n", app.Name, app.ID)
	}
	return nil
}

func GetAppInfoAction(ctx context.Context, cmd *cli.Command) error {
	var app comm.AppView
	resp, err := utils.GetAPIClient().R().SetPathParams(map[string]string{
		"app_id": cmd.Args().First(),
	}).SetResult(&app).Get("/api/v1/applications/{app_id}")
	if err != nil {
		fmt.Printf("Failed to get app info: %s\n", err)
	}

	if resp.StatusCode() != 200 {
		fmt.Printf("Failed to get app info: %s\n", resp.String())
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
	fmt.Printf("Deleting %s...\n", cmd.Args().First())
	return nil
}
