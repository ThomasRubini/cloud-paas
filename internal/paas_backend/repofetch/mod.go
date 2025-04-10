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

// Check if a given directory path exists
func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// called on every repository on a schedule to fetch them, and deploy any needed environments
func FetchAndDeployRepository(state utils.State, project models.DBApplication) error {
	if project.SourceURL == "" {
		logrus.Infof("Skipping %s (empty source URL)", project.Name)
		return nil
	}

	// Init repo if it doesn't exist
	if !isDir(project.GetPath()) {
		err := setupRepo(project, project.GetPath())
		if err != nil {
			return fmt.Errorf("error doing a repository setup : %w", err)
		}
	}

	// Collect branches data before fetching
	oldBranches, err := getAllEnvBranchesLastCommit(project)
	if err != nil {
		return fmt.Errorf("error getting all env branches last commit: %w", err)
	}
	logrus.Debugf("Collected %v branches for project %s before fetching", len(oldBranches), project.Name)

	// Fetch data from remote
	err = fetchRepoChanges(state, project)
	if err != nil {
		return fmt.Errorf("error fetching repository: %w", err)
	}

	// Collect branches data after fetching
	newBranches, err := getAllEnvBranchesLastCommit(project)
	if err != nil {
		return fmt.Errorf("error getting all env branches last commit: %w", err)
	}
	logrus.Debugf("Collected %v branches for project %s after fetching", len(newBranches), project.Name)

	// Check if the commits have changed by comparing branches data of before & after fetching
	for _, env := range project.Envs {
		if oldBranches[env.Branch] != newBranches[env.Branch] {
			logrus.Debugf("New commit for env %v on branch %v", env.Name, env.Branch)
			err := logic.HandleEnvironmentUpdate(state, project, env)
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

	var projects []models.DBApplication
	res := state.Db.Model(&models.DBApplication{}).Preload("Envs").Find(&projects)
	if res.Error != nil {
		return fmt.Errorf("error fetching project names: %w", res.Error)
	}
	logrus.Infof("Found %d projects to fetch", len(projects))

	for _, project := range projects {
		logrus.Debugf("Handling fetching project %v", project.Name)
		err := FetchAndDeployRepository(state, project)
		if err != nil {
			logrus.Errorf("error handling cron update for project %s: %v", project.Name, err)
		}
	}

	return nil
}

// Get the last commit of all branches matching an environment for a given project
// Returns a map of environment name -> last correspondig branch commit hash
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
