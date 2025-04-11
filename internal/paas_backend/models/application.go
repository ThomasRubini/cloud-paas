package models

import (
	"path"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
)

type DBApplication struct {
	BaseModel
	// Immutable ([A-Z][a-z][0-9]-)+
	Name       string `gorm:"unique;not null;default:null"`
	Desc       string
	SourceURL  string
	AutoDeploy bool // TODO move to environment
}

func (DBApplication) TableName() string {
	return "db_applications"
}

func (p DBApplication) GetPath(conf *config.Config) string {
	rootPath := conf.REPO_DIR
	// TODO use project ID for folder
	return path.Join(rootPath, p.Name)
}
