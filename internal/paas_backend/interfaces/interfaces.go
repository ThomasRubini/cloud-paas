package interfaces

import (
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
)

type Logic interface {
	HandleEnvironmentUpdate(app models.DBApplication, env models.DBEnvironment) error
}
