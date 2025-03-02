package models

type DBDeployment struct {
	ParentProject DBApplication
	Domain        string
	Environement  string
}
