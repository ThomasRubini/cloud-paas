package tests

import (
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	// Test health URL
	w := makeRequest(webServer, "GET", "/health", nil)
	assert.Equal(t, 200, w.Code)

	w = makeRequest(webServer, "GET", "/nonexistent", nil)
	assert.Equal(t, 404, w.Code)
}

func TestGetNoApps(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	w := makeRequest(webServer, "GET", "/api/v1/applications", nil)
	assert.Equal(t, 200, w.Code)

	var apps []comm.AppView
	fromJson(w.Body, &apps)
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
	var data map[string]interface{}
	fromJson(w.Body, &data)
	appView.ID = uint(data["id"].(float64))

	// GET requets to check if it was inserted
	w = makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	// Check app content + returned id
	var apps []comm.AppView
	fromJson(w.Body, &apps)
	assert.Equal(t, []comm.AppView{appView}, apps)
}
