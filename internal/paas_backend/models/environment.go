package models

type DBEnvironment struct {
	ParentProject DBApplication
	Domain        string
	Environement  string
}
