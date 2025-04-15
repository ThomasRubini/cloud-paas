package logic

import (
	"encoding/json"
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/deploy"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/imgbuild"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
)

type LogicImpl struct {
	State utils.State
}

// Event called when a repository is updated
func (l LogicImpl) HandleEnvironmentUpdate(app models.DBApplication, env models.DBEnvironment) error {
	// At this point the repository as already been updated

	// Rebuild the image using the updated repository
	imageTag, err := imgbuild.BuildGitBranch(l.State, app, env)
	if err != nil {
		return fmt.Errorf("error building image: %w", err)
	}

	// Upload it to the registry
	err = UploadToRegistry(l.State, imageTag)
	if err != nil {
		return fmt.Errorf("error uploading image to registry: %w", err)
	}

	// Get the exposed port from the image
	port := imgbuild.GetExposedPort(l.State.DockerClient, imageTag)

	parsedEnv := make(map[string]any)
	parsedEnvStr := make(map[string]string)
	if env.EnvVars != "" {
		err = json.Unmarshal([]byte(env.EnvVars), &parsedEnv)
		if err != nil {
			return fmt.Errorf("error unmarshalling environment variables: %w", err)
		}
		parsedEnvStr = make(map[string]string)
		//avoid type misparsing from Kubernetes
		for k, v := range parsedEnv {
			parsedEnvStr[k] = fmt.Sprintf("%v", v)
		}
	}
	// Redeploy to kubernetes using the new image
	err = deploy.DeployEnv(env, deploy.Options{
		ImageTag:    imageTag,
		ExposedPort: *port,
		Namespace:   fmt.Sprintf("%s-%s", config.Get().KUBE_NAMESPACE_PREFIX, app.Name),
		ReleaseName: fmt.Sprintf("%s-%s", app.Name, env.Name),
		EnvVars:     parsedEnvStr,
	})
	if err != nil {
		return fmt.Errorf("error deploying app: %w", err)
	}

	return nil
}

func (l LogicImpl) HandleEnvironmentDeletion(app models.DBApplication, env models.DBEnvironment) error {
	// Delete the environment from kubernetes
	err := deploy.UninstallEnv(env, deploy.Options{
		Namespace:   app.Name,
		ReleaseName: fmt.Sprintf("paas-%s-%s", app.Name, env.Name),
	})
	if err != nil {
		return fmt.Errorf("error deleting environment: %w", err)
	}

	return nil
}
