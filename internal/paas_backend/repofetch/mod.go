package repofetch

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/logic"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
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
func HandleRepository(state utils.State, project models.DBApplication) error {
	if project.SourceURL == "" {
		logrus.Infof("Skipping %s (empty source URL)", project.Name)
	} else {
		commits := make(map[string]string)
		var err error
		if isDir(project.GetPath()) {
			commits, err = getAllEnvBranchesLastCommit(project)
			if err != nil {
				return fmt.Errorf("error getting all env branches last commit: %w", err)
			}
		}
		err = fetchRepository(state, project)
		if err != nil {
			return fmt.Errorf("error fetching repository: %w", err)
		}

		new_commits, err := getAllEnvBranchesLastCommit(project)
		if err != nil {
			return fmt.Errorf("error getting all env branches last commit: %w", err)
		}
		// Check if the commits have changed
		for _, env := range project.Envs {
			if commits[env.Branch] != new_commits[env.Branch] {
				logrus.Info("New commit for env ", env.Name, " on branch ", env.Branch)
				err := logic.HandleEnvironmentUpdate(env)
				if err != nil {
					return fmt.Errorf("error handling repository update: %w", err)
				}
			}
		}
	}

	return nil
}

func handleRepositories() error {
	logrus.Info("Pulling repositories at", time.Now())

	// Get state
	state := utils.GetState()

	var projects []models.DBApplication
	res := state.Db.Model(&models.DBApplication{}).Preload("Envs").Find(&projects)
	if res.Error != nil {
		return fmt.Errorf("error fetching project names: %w", res.Error)
	}
	logrus.Infof("Found %d projects to pull", len(projects))

	for _, project := range projects {
		err := HandleRepository(state, project)
		if err != nil {
			logrus.Errorf("error handling cron update for project %s: %v", project.Name, err)
		}
	}

	return nil
}

func getAllEnvBranchesLastCommit(project models.DBApplication) (map[string]string, error) {
	dir := project.GetPath()
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("error opening repository for project %v : %w", project, err)
	}
	// Get the remote branches
	refIter, err := repo.Storer.IterReferences()
	if err != nil {
		return nil, err
	}
	branches := storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsRemote()
		}, refIter)

	branchesLastCommit := make(map[string]string)
	//TODO : Optimize this shit
	for _, env := range project.Envs {
		err = branches.ForEach(func(branch *plumbing.Reference) error {
			if strings.TrimPrefix(branch.Name().String(), "refs/remotes/origin/") == env.Branch {
				branchesLastCommit[env.Name] = branch.Hash().String()
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error iterating branches for project %v: %w", project, err)
		}
	}
	return branchesLastCommit, nil
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
