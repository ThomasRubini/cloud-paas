package repofetch

import (
	"errors"
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

func getAuth(project models.DBApplication) transport.AuthMethod {
	if project.SourcePassword != "" {
		return &http.BasicAuth{
			Username: project.SourceUsername,
			Password: project.SourcePassword,
		}
	}
	return nil
}

func initRepoIfNotExists(project models.DBApplication, dir string) error {
	repo, err := git.PlainInit(dir, false)
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

func fetchRepoChanges(project models.DBApplication, dir string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("error opening repository: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %v", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       getAuth(project),
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("error fetching repository: %v", err)
	}
	return nil
}

func pullRepository(project models.DBApplication) error {
	dir := project.GetPath()

	logrus.Debugf("Pulling repository %v at %v", project.Name, dir)

	if !isDir(dir) {
		err := initRepoIfNotExists(project, dir)
		if err != nil {
			return err
		}
	}

	err := fetchRepoChanges(project, dir)
	if err != nil {
		return err
	}

	return nil
}
