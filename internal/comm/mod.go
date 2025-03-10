// Structs common to frontend and backend
package comm

type IdResponse struct {
	ID uint `json:"id"`
}

type CreateAppRequest struct {
	Name           string
	Desc           string
	SourceURL      string
	SourceUsername string
	SourcePassword string
	AutoDeploy     bool
}

type AppView struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	SourceURL  string `json:"source_url"`
	AutoDeploy bool   `json:"auto_deploy"`
}

type EnvView struct {
	ID     uint `json:"id"`
	Domain string
	Name   string
}

type CreateEnvRequest struct {
	Domain string
	Name   string
}
