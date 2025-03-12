package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func splitImage(image string) (string, string) {
	split := strings.Split(image, ":")
	if len(split) == 1 {
		return split[0], "latest"
	}
	return split[0], split[1]
}

func getDeploymentName(env models.DBEnvironment) string {
	return fmt.Sprintf("%v-%v", env.Application.Name, env.Name)
}

func generateChart(env models.DBEnvironment, options Options) (*chart.Chart, error) {
	deploymentPath := filepath.Join("dist", "deployment.yaml")
	deploymentData, err := os.ReadFile(deploymentPath)
	if err != nil {
		return nil, err
	}

	imageRepo, imageTag := splitImage(options.ImageTag)

	myChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:       getDeploymentName(env),
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
			"namespace":    options.Namespace,
			"replicaCount": 1,
			"image": map[string]interface{}{
				"repository": imageRepo,
				"tag":        imageTag,
			},
			"service": map[string]interface{}{
				"port": options.ExposedPort,
			},
		},
	}

	return myChart, nil
}

func installHelmChart(myChart *chart.Chart, env models.DBEnvironment, options Options) (*release.Release, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), options.Namespace, "memory", nil); err != nil {
		return nil, err
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = env.Application.Name
	install.Namespace = options.Namespace
	install.Wait = true
	install.Atomic = true

	return install.Run(myChart, myChart.Values)
}

/*
Inputs:
- Deployement (in the cloud-pass sense, not the kubernetes sense)
  - Name
  - Image
  - TCP Port exposed inside the image
*/

type Options struct {
	Namespace   string
	ImageTag    string
	ExposedPort int
}

func DeployApp(env models.DBEnvironment, options Options) error {
	logrus.Debugf("Deploying app %v:%v", env.Application.Name, env.Name)

	myChart, err := generateChart(env, options)
	if err != nil {
		return fmt.Errorf("error generating chart: %v", err)
	}

	_, err = installHelmChart(myChart, env, options)
	if err != nil {
		return fmt.Errorf("error installing chart: %v", err)
	}

	logrus.Debugf("Deployed app %v:%v successfully", env.Application.Name, env.Name)
	return nil
}
