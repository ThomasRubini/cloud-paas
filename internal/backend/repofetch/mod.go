package repofetch

import (
	"cloud-paas/internal/backend/config"
	"cloud-paas/internal/backend/models"
	"cloud-paas/internal/backend/state"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/sirupsen/logrus"
)

func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func handleRepository(project models.DBProject) error {
	err := pullRepository(project)
	if err != nil {
		return fmt.Errorf("error pulling repository: %v", err)
	}

	return nil
}

func pullRepository(project models.DBProject) error {
	p := config.Get().REPO_DIR
	logrus.Debugf("Pulling repository %v at %v", project.Name, p)

	// TODO use project ID for cloning
	dir := path.Join(p, project.Name)
	if !isDir(dir) {
		repo, err := git.PlainInit(dir, true)
		if err != nil {
			return fmt.Errorf("error initializing repository: %v", err)
		}

		_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
			Name: "origin",
			URLs: []string{project.SourceURL},
		})
		if err != nil {
			return fmt.Errorf("error adding remote: %v", err)
		}
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("error opening repository: %v", err)
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("error fetching repository: %v", err)
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
			logrus.Errorf("Error pulling repository for project %v: %v", project, err)
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
