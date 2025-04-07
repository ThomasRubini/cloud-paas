package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

func generateChart(env models.DBEnvironment, options Options) (*chart.Chart, error) {
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

func installHelmChart(helmConfig *action.Configuration, myChart *chart.Chart, options Options) (*release.Release, error) {
	install := action.NewInstall(helmConfig)
	install.ReleaseName = options.ReleaseName
	install.Namespace = options.Namespace
	install.Wait = true
	install.Atomic = true
	install.CreateNamespace = true
	install.Timeout = 30 * time.Second

	resp, err := install.Run(myChart, myChart.Values)
	if err != nil {
		return nil, fmt.Errorf("error running install: %w", err)
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
	logrus.Debugf("Deploying app %v:%v", env.Application.Name, env.Name)

	myChart, err := generateChart(env, options)
	if err != nil {
		return fmt.Errorf("error generating chart: %w", err)
	}

	_, err = installHelmChart(helmConfig, myChart, options)
	if err != nil {
		return fmt.Errorf("error installing chart: %w", err)
	}

	logrus.Debugf("Deployed app %v:%v successfully", env.Application.Name, env.Name)
	return nil
}
