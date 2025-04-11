package deploy

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func GetHelmConfig(namespace string) (*action.Configuration, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", log.Printf); err != nil {
		return nil, fmt.Errorf("error initializing config: %w", err)
	}
	return actionConfig, nil
}

func generateChart(options Options, env models.DBEnvironment) (*chart.Chart, error) {

	templates := []*chart.File{}
	err := filepath.Walk("assets/helm_templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			relativePath, relErr := filepath.Rel("assets/helm_templates", path)
			if relErr != nil {
				return relErr
			}
			templates = append(templates, &chart.File{
				Name: filepath.Join("templates", relativePath),
				Data: data,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error reading helm templates: %w", err)
	}

	fmt.Println("AAAAAAAAAa")
	fmt.Printf("%v\n", options)

	myChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:       options.ReleaseName,
			APIVersion: "v2",
			Version:    "0.1.0",
		},
		Templates: templates,
		Values: map[string]interface{}{
			"deploymentName": options.ReleaseName,
			"namespace":      options.Namespace,
			"replicaCount":   1, // TODO: Add autoscaling later
			"image":          options.ImageTag,
			"containerPort":  options.ExposedPort,
			"domain":         env.Domain,
			"envVariables":   options.EnvVars,
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
		return installHelmChart(helmConfig, myChart, options, isReleaseUninstalled(versions))
	} else {
		return upgradeHelmChart(helmConfig, myChart, options)
	}
}

func installHelmChart(helmConfig *action.Configuration, myChart *chart.Chart, options Options, isUninstalled bool) (*release.Release, error) {
	logrus.Debugf("using Install action for release %v (isUninstalled=%v)", options.ReleaseName, isUninstalled)
	install := action.NewInstall(helmConfig)
	install.ReleaseName = options.ReleaseName
	install.CreateNamespace = true
	install.Replace = isUninstalled

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
	logrus.Debugf("Using Upgrade action for release %v", options.ReleaseName)
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
	EnvVars     map[string]string
}

func DeployEnv(env models.DBEnvironment, options Options) error {
	logrus.Debugf("Deploying release %v", options.ReleaseName)
	helmConfig, err := GetHelmConfig(options.Namespace)
	if err != nil {
		return fmt.Errorf("error getting helm config: %w", err)
	}

	myChart, err := generateChart(options, env)
	if err != nil {
		return fmt.Errorf("error generating chart: %w", err)
	}

	_, err = applyHelmChart(helmConfig, myChart, options)
	if err != nil {
		return fmt.Errorf("error installing chart: %w", err)
	}

	logrus.Debugf("Deployed env %v successfully", options.ReleaseName)
	return nil
}

func UninstallEnv(env models.DBEnvironment, options Options) error {
	logrus.Debugf("Uninstalling release %v", options.ReleaseName)
	helmConfig, err := GetHelmConfig(options.Namespace)
	if err != nil {
		return fmt.Errorf("error getting helm config: %w", err)
	}

	uninstall := action.NewUninstall(helmConfig)
	uninstall.Timeout = 30 * time.Second
	uninstall.Wait = true

	_, err = uninstall.Run(options.ReleaseName)
	if err != nil {
		return fmt.Errorf("error uninstalling chart: %w", err)
	}

	logrus.Debugf("Uninstalled env %v successfully", options.ReleaseName)
	return nil
}
