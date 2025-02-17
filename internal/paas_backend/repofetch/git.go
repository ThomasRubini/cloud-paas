package repofetch

import (
	"errors"
	"fmt"
	"path"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/sirupsen/logrus"
)

func initRepoIfNotExists(project models.DBProject, dir string) error {
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
	return nil
}

func fetchRepoChanges(dir string) error {
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

func pullRepository(project models.DBProject) error {
	p := config.Get().REPO_DIR
	// TODO use project ID for folder
	dir := path.Join(p, project.Name)

	logrus.Debugf("Pulling repository %v at %v", project.Name, p)

	if !isDir(dir) {
		err := initRepoIfNotExists(project, dir)
		if err != nil {
			return err
		}
	}

	err := fetchRepoChanges(dir)
	if err != nil {
		return err
	}

	return nil
}
