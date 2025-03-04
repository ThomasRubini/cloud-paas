package models

import (
	"fmt"
	"path"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/state"
	"gorm.io/gorm"
)

type DBApplication struct {
	gorm.Model
	// Immutable ([A-Z][a-z][0-9]-)+
	Name       string
	Desc       string
	SourceURL  string
	AutoDeploy bool
}

func (p DBApplication) GetPath() string {
	rootPath := config.Get().REPO_DIR
	// TODO use project ID for folder
	return path.Join(rootPath, p.Name)
}

func (p DBApplication) SetSourceCredentials(username, password string) error {
	sp := state.Get().SecretsProvider

	if username == "" && password == "" {
		// Delete
		if err := sp.DeleteSecret(fmt.Sprintf("%v.username", p.ID)); err != nil {
			return fmt.Errorf("could not delete source username: %v", err)
		}

		if err := sp.DeleteSecret(fmt.Sprintf("%v.password", p.ID)); err != nil {
			return fmt.Errorf("could not delete source password: %v", err)
		}
	} else {
		// Set
		if err := sp.SetSecret(fmt.Sprintf("%v.username", p.ID), username); err != nil {
			return fmt.Errorf("could not set source username: %v", err)
		}

		if err := sp.SetSecret(fmt.Sprintf("%v.password", p.ID), password); err != nil {
			return fmt.Errorf("could not set source password: %v", err)
		}
	}

	return nil
}

func (p DBApplication) GetSourceCredentials() (string, string, error) {
	sp := state.Get().SecretsProvider

	username, err := sp.GetSecret(fmt.Sprintf("%v.username", p.ID))
	if err != nil {
		return "", "", fmt.Errorf("could not get source username: %v", err)
	}

	password, err := sp.GetSecret(fmt.Sprintf("%v.password", p.ID))
	if err != nil {
		return "", "", fmt.Errorf("could not get source password: %v", err)
	}

	return username, password, nil
}
