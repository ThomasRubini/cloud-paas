package models

type DBDeployment struct {
	ParentProject DBProject
	Environement  string
}
