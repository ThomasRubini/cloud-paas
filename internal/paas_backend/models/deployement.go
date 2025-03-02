package models

type DBDeployment struct {
	ParentProject DBProject
	Domain        string
	Environement  string
}
