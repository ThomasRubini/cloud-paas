package tests

import (
	"testing"

	"fmt"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	// Test health URL
	w := makeOKRequest(t, webServer, "GET", "/health", nil)
	assert.Equal(t, "OK", toString(w.Body))

	w = makeRequest(webServer, "GET", "/nonexistent", nil)
	assert.Equal(t, 404, w.Code)
}

func TestGetNoApps(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	w := makeRequest(webServer, "GET", "/api/v1/applications", nil)
	assert.Equal(t, 200, w.Code)

	var apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, 0, len(apps))
}

func TestOneApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name: "test",
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Make POST request
	w := makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))
	data := fromJson[map[string]uint](w.Body)
	appView.ID = data["id"]

	// GET requets to check if it was inserted
	w = makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	// Check app content + returned id
	var apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, []comm.AppView{appView}, apps)
}

func TestDeleteApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name: "test",
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Insert app
	makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))

	// Verify app was inserted
	w := makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)
	apps := fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, 1, len(apps))

	// Delete app
	makeOKRequest(t, webServer, "DELETE", fmt.Sprintf("/api/v1/applications/%v", apps[0].ID), nil)

	// GET requets to check if it was deleted
	w = makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	// Check app content + returned id
	apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, 0, len(apps))
}
