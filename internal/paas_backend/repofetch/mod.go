package repofetch

import (
	"fmt"
	"os"
	"time"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/logic"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
		logrus.Debug("Skipping project with empty source URL")
	}

	commits := getAllEnvBranchesLastCommit(project)

	err := fetchRepository(state, project)
	if err != nil {
		return fmt.Errorf("error fetching repository: %v", err)
	}

	new_commits := getAllEnvBranchesLastCommit(project)

	for _, env := range project.Envs {
		if commits[env.Branch] != new_commits[env.Branch] {
			err := logic.HandleEnvironmentUpdate(env)
			if err != nil {
				return fmt.Errorf("error handling repository update: %v", err)
			}
		}
	}

	for _, env := range project.Envs {
		err = logic.HandleEnvironmentUpdate(env)
		if err != nil {
			return fmt.Errorf("error pulling environment: %v", err)
		}
	}

	if err != nil {
		return fmt.Errorf("error handling repository update: %v", err)
	}

	return nil
}

func handleRepositories() error {
	logrus.Info("Pulling repositories at", time.Now())

	// Get state
	state := utils.GetState()

	var projects []models.DBApplication
	res := state.Db.Model(&models.DBApplication{}).Find(&projects)
	if res.Error != nil {
		return fmt.Errorf("error fetching project names: %v", res.Error)
	}
	logrus.Infof("Found %d projects to pull", len(projects))

	for _, project := range projects {
		err := HandleRepository(state, project)
		if err != nil {
			logrus.Errorf("Error handling cron update for project %v: %v", project, err)
		}
	}

	return nil
}

func getAllEnvBranchesLastCommit(project models.DBApplication) map[string]string {
	dir := project.GetPath()
	repo, err := git.PlainOpen(dir)
	if err != nil {
		logrus.Errorf("Error opening repository for project %v : %v", project, err)
	}
	branches, err := repo.Branches()
	if err != nil {
		logrus.Errorf("Error getting all branches for project %v : %v", project, err)
	}
	branchesLastCommit := make(map[string]string)
	//TODO : Optimize this shit
	for _, env := range project.Envs {
		branches.ForEach(func(branch *plumbing.Reference) error {
			if branch.Name().String() == env.Branch {
				branchesLastCommit[env.Name] = branch.Hash().String()
			}
			return nil
		})
	}
	return branchesLastCommit
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
