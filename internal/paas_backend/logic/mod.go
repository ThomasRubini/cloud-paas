package logic

import (
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/imgbuild"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
)

// Event called when a repository is updated
func HandleRepositoryUpdate(project models.DBApplication) error {
	// At this point the repository as already been updated

	// Rebuild the image using the updated repository
	// TODO what to name the tags ?
	err := imgbuild.Build(project.GetPath(), project.Name)
	if err != nil {
		return fmt.Errorf("error building image: %v", err)
	}

	// Redeploy using the new image
	panic("Thomas appelle ça")
	/*err = deploy.DeployApp(project)
	if err != nil {
		return fmt.Errorf("error deploying app: %v", err)
	}
	return nil
	*/
}
