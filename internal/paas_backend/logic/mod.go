package logic

import (
	"fmt"

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
	imageTag := fmt.Sprintf("%s:%s", project.Name, env.Name)
	err := imgbuild.Build(project.GetPath(), imageTag)
	if err != nil {
		return fmt.Errorf("error building image: %v", err)
	}

	port := imgbuild.GetExposedPort(imageTag)

	// Redeploy using the new image
	err = deploy.DeployApp(env, imageTag, env.Application.Name, *port)
	if err != nil {
		return fmt.Errorf("error deploying app: %v", err)
	}

	return nil
}
