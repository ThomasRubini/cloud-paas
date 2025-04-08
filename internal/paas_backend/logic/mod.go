package logic

import (
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/deploy"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/imgbuild"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
)

// Event called when a repository is updated
func HandleEnvironmentUpdate(state utils.State, app models.DBApplication, env models.DBEnvironment) error {
	// At this point the repository as already been updated

	// Rebuild the image using the updated repository
	// TODO what to name the tags ?
	imageTag := fmt.Sprintf("%s/%s:%s", config.Get().REGISTRY_REPO_URI, app.Name, env.Name)
	err := imgbuild.BuildGitBranch(state.DockerClient, app.GetPath(), env.Branch, imageTag)
	if err != nil {
		return fmt.Errorf("error building image: %w", err)
	}

	err = UploadToRegistry(state.DockerClient, imageTag)
	if err != nil {
		return fmt.Errorf("error uploading image to registry: %w", err)
	}

	port := imgbuild.GetExposedPort(state.DockerClient, imageTag)

	// Redeploy using the new image
	err = deploy.DeployApp(state.HelmConfig, env, deploy.Options{
		ImageTag:    imageTag,
		ExposedPort: *port,
		Namespace:   app.Name,
		ReleaseName: fmt.Sprintf("paas-%s-%s", app.Name, env.Name),
	})
	if err != nil {
		return fmt.Errorf("error deploying app: %w", err)
	}

	return nil
}
