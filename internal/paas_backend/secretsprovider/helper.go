package secretsprovider

import (
	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
)

type Helper struct {
	Core
}

func (h Helper) SetSourceCredentials(app models.DBApplication, username, password string) error {
	if username == "" && password == "" {
		// Delete
		if err := h.DeleteSecret(fmt.Sprintf("%v.username", app.ID)); err != nil {
			return fmt.Errorf("could not delete source username: %w", err)
		}

		if err := h.DeleteSecret(fmt.Sprintf("%v.password", app.ID)); err != nil {
			return fmt.Errorf("could not delete source password: %w", err)
		}
	} else {
		// Set
		if err := h.SetSecret(fmt.Sprintf("%v.username", app.ID), username); err != nil {
			return fmt.Errorf("could not set source username: %w", err)
		}

		if err := h.SetSecret(fmt.Sprintf("%v.password", app.ID), password); err != nil {
			return fmt.Errorf("could not set source password: %w", err)
		}
	}

	return nil
}

func (h Helper) GetSourceCredentials(app models.DBApplication) (string, string, error) {
	username, err := h.GetSecret(fmt.Sprintf("%v.username", app.ID))
	if err != nil {
		return "", "", fmt.Errorf("could not get source username: %w", err)
	}

	password, err := h.GetSecret(fmt.Sprintf("%v.password", app.ID))
	if err != nil {
		return "", "", fmt.Errorf("could not get source password: %w", err)
	}

	return username, password, nil
}
