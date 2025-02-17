package models

import (
	"path"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
)

type DBProject struct {
	Name           string
	Desc           string
	SourceURL      string
	SourceUsername string
	SourcePassword string
	AutoDeploy     bool
}

func (p DBProject) GetPath() string {
	rootPath := config.Get().REPO_DIR
	// TODO use project ID for folder
	return path.Join(rootPath, p.Name)
}
