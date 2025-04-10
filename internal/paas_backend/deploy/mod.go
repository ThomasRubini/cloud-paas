package deploy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func generateChart(options Options) (*chart.Chart, error) {
	deploymentPath := filepath.Join("assets", "app_deployments", "deployment.yaml")
	deploymentData, err := os.ReadFile(deploymentPath)
	if err != nil {
		return nil, err
	}

	myChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:       options.ReleaseName,
			APIVersion: "v2",
			Version:    "0.1.0",
		},
		Templates: []*chart.File{
			{
				Name: "templates/deployment.yaml",
				Data: deploymentData,
			},
		},
		Values: map[string]interface{}{
			"deploymentName": options.ReleaseName,
			"namespace":      options.Namespace,
			"replicaCount":   1,
			"image":          options.ImageTag,
			"containerPort":  options.ExposedPort,
		},
	}

	return myChart, nil
}

func isReleaseUninstalled(versions []*release.Release) bool {
	return len(versions) > 0 && versions[len(versions)-1].Info.Status == release.StatusUninstalled
}

// Install or Upgrade the Helm chart
// Credits: https://github.com/helm/helm/blob/d8ca55fc669645c10c0681d49723f4bb8c0b1ce7/pkg/cmd/upgrade.go
func applyHelmChart(helmConfig *action.Configuration, myChart *chart.Chart, options Options) (*release.Release, error) {
	histClient := action.NewHistory(helmConfig)
	histClient.Max = 1
	versions, err := histClient.Run(options.ReleaseName)
	if errors.Is(err, driver.ErrReleaseNotFound) || isReleaseUninstalled(versions) {
		return installHelmChart(helmConfig, myChart, options)
	} else {
		return upgradeHelmChart(helmConfig, myChart, options)
	}
}

func installHelmChart(helmConfig *action.Configuration, myChart *chart.Chart, options Options) (*release.Release, error) {
	install := action.NewInstall(helmConfig)
	install.ReleaseName = options.ReleaseName
	install.CreateNamespace = true

	install.Namespace = options.Namespace
	install.Wait = true
	install.Atomic = true
	install.Timeout = 30 * time.Second

	resp, err := install.Run(myChart, myChart.Values)
	if err != nil {
		return nil, fmt.Errorf("error running install: %w", err)
	}

	return resp, nil
}

func upgradeHelmChart(helmConfig *action.Configuration, myChart *chart.Chart, options Options) (*release.Release, error) {
	upgrade := action.NewUpgrade(helmConfig)

	upgrade.Namespace = options.Namespace
	upgrade.Wait = true
	upgrade.Atomic = true
	upgrade.Timeout = 30 * time.Second

	resp, err := upgrade.Run(options.ReleaseName, myChart, myChart.Values)
	if err != nil {
		return nil, fmt.Errorf("error running upgrade: %w", err)
	}

	return resp, nil
}

type Options struct {
	Namespace   string
	ImageTag    string
	ExposedPort int
	ReleaseName string
}

func DeployApp(helmConfig *action.Configuration, env models.DBEnvironment, options Options) error {
	logrus.Debugf("Deploying release %v", options.ReleaseName)

	myChart, err := generateChart(options)
	if err != nil {
		return fmt.Errorf("error generating chart: %w", err)
	}

	_, err = applyHelmChart(helmConfig, myChart, options)
	if err != nil {
		return fmt.Errorf("error installing chart: %w", err)
	}

	logrus.Debugf("Deployed app %v successfully", options.ReleaseName)
	return nil
}
