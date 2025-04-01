package logic

import (
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/deploy"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/imgbuild"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
)

// Event called when a repository is updated
func HandleEnvironmentUpdate(env models.DBEnvironment) error {
	// At this point the repository as already been updated

	// Rebuild the image using the updated repository
	// TODO what to name the tags ?
	project := env.Application

	imageTag := fmt.Sprintf("%s/%s:%s", config.Get().REGISTRY_REPO_URI, project.Name, env.Name)
	err := imgbuild.Build(project.GetPath(), imageTag)
	if err != nil {
		return fmt.Errorf("error building image: %w", err)
	}

	err = UploadToRegistry(imageTag)
	if err != nil {
		return fmt.Errorf("error uploading image to registry: %w", err)
	}

	port := imgbuild.GetExposedPort(imageTag)

	// Redeploy using the new image
	err = deploy.DeployApp(env, deploy.Options{
		ImageTag:    imageTag,
		ExposedPort: *port,
		Namespace:   env.Application.Name,
	})
	if err != nil {
		return fmt.Errorf("error deploying app: %w", err)
	}

	return nil
}
