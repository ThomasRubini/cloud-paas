package models

import (
	"path"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
)

type DBApplication struct {
	// Immutable ([A-Z][a-z][0-9]-)+
	Name           string
	Desc           string
	SourceURL      string
	SourceUsername string
	SourcePassword string
	AutoDeploy     bool
}

func (p DBApplication) GetPath() string {
	rootPath := config.Get().REPO_DIR
	// TODO use project ID for folder
	return path.Join(rootPath, p.Name)
}
