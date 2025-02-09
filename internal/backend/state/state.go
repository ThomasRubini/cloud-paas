package state

import (
	"gorm.io/gorm"
)

type T struct {
	Db *gorm.DB
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
