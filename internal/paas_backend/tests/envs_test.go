package tests

import (
	"fmt"
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetNoEnvs(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	newApp := comm.CreateAppRequest{
		Name: "test1",
	}
	makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(newApp))

	w := makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), nil)

	var envs = fromJson[[]comm.EnvView](w.Body)
	assert.Equal(t, 0, len(envs))
}

func TestGetEnv(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	newApp := comm.CreateAppRequest{
		Name: "test1",
	}
	makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(newApp))

	newEnv := comm.CreateEnvRequest{
		Name: "test1",
	}
	var finalEnv comm.EnvView
	utils.CopyFields(&newEnv, &finalEnv)
	w := makeOKRequest(t, webServer, "POST", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), toJson(newEnv))
	finalEnv.ID = fromJson[comm.EnvView](w.Body).ID

	w = makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), nil)

	var envs = fromJson[[]comm.EnvView](w.Body)
	assert.Equal(t, []comm.EnvView{finalEnv}, envs)
}
