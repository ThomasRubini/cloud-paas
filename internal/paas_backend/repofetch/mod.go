package repofetch

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/sirupsen/logrus"
)

// Check if a given directory path exists
func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// called on every repository on a schedule to fetch them, and deploy any needed environments
func FetchAndDeployRepository(state utils.State, app models.AppWithEnvs) error {
	if app.SourceURL == "" {
		logrus.Infof("Skipping %s (empty source URL)", app.Name)
		return nil
	}

	// Init repo if it doesn't exist
	repoPath := app.GetPath(state.Config)
	if !isDir(repoPath) {
		err := setupRepo(app.DBApplication, repoPath)
		if err != nil {
			return fmt.Errorf("error doing a repository setup : %w", err)
		}
	}

	// Collect branches data before fetching
	oldBranches, err := getAllEnvBranchesLastCommit(state, app)
	if err != nil {
		return fmt.Errorf("error getting all env branches last commit: %w", err)
	}
	logrus.Debugf("Collected %v branches for project %s before fetching", len(oldBranches), app.Name)

	// Fetch data from remote
	err = FetchRepoChanges(state, app.DBApplication)
	if err != nil {
		return fmt.Errorf("error fetching repository: %w", err)
	}

	// Collect branches data after fetching
	newBranches, err := getAllEnvBranchesLastCommit(state, app)
	if err != nil {
		return fmt.Errorf("error getting all env branches last commit: %w", err)
	}
	logrus.Debugf("Collected %v branches for project %s after fetching", len(newBranches), app.Name)

	// Check if the commits have changed by comparing branches data of before & after fetching
	for _, env := range app.Envs {
		logrus.Debugf("Comparing env %s: %s -> %s", env.Name, oldBranches[env.Name], newBranches[env.Name])
		if oldBranches[env.Name] != newBranches[env.Name] {
			logrus.Debugf("New commit for env %v on branch %v", env.Name, env.Branch)
			err := state.LogicModule.HandleEnvironmentUpdate(app.DBApplication, env)
			if err != nil {
				return fmt.Errorf("error handling repository update: %w", err)
			}
		}
	}

	return nil
}

// Fetch every repository, and deploy needed environments
func handleRepositories() error {
	logrus.Info("Fetching repositories due to recurring task")

	// Get state
	state := utils.GetState()

	var apps []models.AppWithEnvs
	res := state.Db.Model(&models.AppWithEnvs{}).Preload("Envs").Find(&apps)
	if res.Error != nil {
		return fmt.Errorf("error fetching project names: %w", res.Error)
	}
	logrus.Infof("Found %d projects to fetch", len(apps))

	for _, app := range apps {
		logrus.Debugf("Handling fetching project %v", app.Name)
		err := FetchAndDeployRepository(state, app)
		if err != nil {
			logrus.Errorf("error handling cron update for project %s: %v", app.Name, err)
		}
	}

	return nil
}

// Get the last commit of all branches matching an environment for a given project
// Returns a map of environment name -> last correspondig branch commit hash
func getAllEnvBranchesLastCommit(state utils.State, app models.AppWithEnvs) (map[string]string, error) {
	dir := app.GetPath(state.Config)
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("error opening repository for project %v : %w", app.Name, err)
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
	for _, env := range app.Envs {
		err = branches.ForEach(func(branch *plumbing.Reference) error {
			if strings.TrimPrefix(branch.Name().String(), "refs/remotes/origin/") == env.Branch {
				branchesLastCommit[env.Name] = branch.Hash().String()
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error iterating branches for project %v: %w", app.Name, err)
		}
	}
	return branchesLastCommit, nil
}

// Start the repository fetcher task
func Init(period int) {
	go func() {
		for {
			err := handleRepositories()
			if err != nil {
				logrus.Errorf("Error fetching repositories: %v", err)
			}
			logrus.Debugf("Finished fetching repositories, sleeping for %d seconds", period)

			time.Sleep(time.Duration(period) * time.Second)
		}
	}()
}
