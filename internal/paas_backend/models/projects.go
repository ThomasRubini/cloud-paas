package models

type DBProject struct {
	Name           string
	Desc           string
	SourceURL      string
	SourceUsername string
	SourcePassword string
	AutoDeploy     bool
}
