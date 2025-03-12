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
		Name: "prod",
	}
	var finalEnv comm.EnvView
	utils.CopyFields(&newEnv, &finalEnv)
	w := makeOKRequest(t, webServer, "POST", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), toJson(newEnv))
	finalEnv.ID = fromJson[comm.EnvView](w.Body).ID

	w = makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), nil)

	var envs = fromJson[[]comm.EnvView](w.Body)
	assert.Equal(t, []comm.EnvView{finalEnv}, envs)
}

func TestDeleteEnv(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	newApp := comm.CreateAppRequest{
		Name: "test1",
	}
	// Add app
	makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(newApp))

	// Add env
	newEnv := comm.CreateEnvRequest{
		Name: "prod",
	}
	makeOKRequest(t, webServer, "POST", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), toJson(newEnv))

	// get and check that present
	w := makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), nil)
	envs := fromJson[[]comm.EnvView](w.Body)
	assert.Equal(t, 1, len(envs))
	fmt.Printf("%v\n", envs)

	// Delete
	makeOKRequest(t, webServer, "DELETE", fmt.Sprintf("/api/v1/applications/%v/environments/%v", newApp.Name, newEnv.Name), nil)

	// Check that not here anymore
	w = makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v/environments", newApp.Name), nil)
	envs = fromJson[[]comm.EnvView](w.Body)
	assert.Equal(t, 0, len(envs))
}

func TestDeleteNonExistingEnv(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	newApp := comm.CreateAppRequest{
		Name: "test1",
	}
	// Add app
	makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(newApp))

	// Delete
	w := makeRequest(webServer, "DELETE", fmt.Sprintf("/api/v1/applications/%v/environments/nonexisting", newApp.Name), nil)
	assert.Equal(t, 404, w.Code)
}
