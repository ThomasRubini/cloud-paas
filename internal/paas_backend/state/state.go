package state

import (
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/secretsprovider"
	"gorm.io/gorm"
)

type T struct {
	Db              *gorm.DB
	SecretsProvider secretsprovider.SecretsProvider
}

var state *T

func Get() T {
	if state == nil {
		panic("state not set")
	}
	return *state
}

func Set(s T) {
	if state != nil {
		panic("state already set")
	}
	state = &s
}
