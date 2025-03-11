package repofetch

import (
	"errors"
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
)

func getAuth(state utils.State, project models.DBApplication) (transport.AuthMethod, error) {
	username, password, err := state.SecretsProvider.GetSourceCredentials(project)
	if err != nil {
		return nil, fmt.Errorf("error getting source credentials: %v", err)
	}

	return &http.BasicAuth{
		Username: username,
		Password: password,
	}, nil
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

func fetchRepoChanges(state utils.State, project models.DBApplication, dir string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("error opening repository: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %v", err)
	}

	auth, err := getAuth(state, project)
	if err != nil {
		return fmt.Errorf("error getting auth: %v", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("error fetching repository: %v", err)
	}
	return nil
}

func pullRepository(state utils.State, project models.DBApplication) error {
	dir := project.GetPath()

	logrus.Debugf("Pulling repository %v at %v", project.Name, dir)

	if !isDir(dir) {
		err := initRepoIfNotExists(project, dir)
		if err != nil {
			return err
		}
	}

	err := fetchRepoChanges(state, project, dir)
	if err != nil {
		return err
	}

	return nil
}
