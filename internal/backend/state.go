package backend

import (
	"gorm.io/gorm"
)

type BackendState struct {
	db *gorm.DB
}

var state *BackendState

func GetState() BackendState {
	if state == nil {
		panic("state not set")
	}
	return *state
}

func SetState(s BackendState) {
	if state != nil {
		panic("state already set")
	}
	state = &s
}
