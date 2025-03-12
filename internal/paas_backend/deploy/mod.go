package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
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

func generateChart(env models.DBEnvironment, appImage string, appPort int) (*chart.Chart, error) {
	deploymentPath := filepath.Join("dist", "deployment.yaml")
	deploymentData, err := os.ReadFile(deploymentPath)
	if err != nil {
		return nil, err
	}

	imageRepo, imageTag := splitImage(appImage)

	myChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:       env.Application.Name,
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
			"replicaCount": 1,
			"image": map[string]interface{}{
				"repository": imageRepo,
				"tag":        imageTag,
			},
			"service": map[string]interface{}{
				"port": appPort,
			},
		},
	}

	return myChart, nil
}

func installHelmChart(myChart *chart.Chart, env models.DBEnvironment, namespace string) (*release.Release, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "memory", nil); err != nil {
		return nil, err
	}

	install := action.NewInstall(actionConfig)
	install.ReleaseName = env.Application.Name
	install.Namespace = namespace
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
func DeployApp(env models.DBEnvironment, imageTag string, namespace string, exposedPort int) error {

	myChart, err := generateChart(env, imageTag, exposedPort)
	if err != nil {
		return fmt.Errorf("error generating chart: %v", err)
	}

	_, err = installHelmChart(myChart, env, namespace)
	if err != nil {
		return fmt.Errorf("error installing chart: %v", err)
	}

	return nil
}
