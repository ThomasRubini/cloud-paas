package logic

import (
	"fmt"

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

	// Redeploy to kubernetes using the new image
	err = deploy.DeployApp(l.State.HelmConfig, env, deploy.Options{
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
