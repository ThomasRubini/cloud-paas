package repofetch

import (
	"fmt"
	"os"
	"time"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/logic"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/state"
	"github.com/sirupsen/logrus"
)

func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// called on every repository on a schedule to pull them and update them
func handleRepository(project models.DBProject) error {
	err := pullRepository(project)
	if err != nil {
		return fmt.Errorf("error pulling repository: %v", err)
	}

	err = logic.HandleRepositoryUpdate(project)
	if err != nil {
		return fmt.Errorf("error handling repository update: %v", err)
	}

	return nil
}

func handleRepositories() error {
	logrus.Info("Pulling repositories at", time.Now())

	var projects []models.DBProject
	res := state.Get().Db.Model(&models.DBProject{}).Find(&projects)
	if res.Error != nil {
		return fmt.Errorf("error fetching project names: %v", res.Error)
	}
	logrus.Infof("Found %d projects to pull", len(projects))

	for _, project := range projects {
		err := handleRepository(project)
		if err != nil {
			logrus.Errorf("Error handling cron update for project %v: %v", project, err)
		}
	}

	return nil
}

func Init(period int) {
	go func() {
		err := handleRepositories()
		if err != nil {
			logrus.Errorf("Error pulling repositories: %v", err)
		}

		time.Sleep(time.Duration(period) * time.Second)
	}()
}
