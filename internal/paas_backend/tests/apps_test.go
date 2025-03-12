package tests

import (
	"fmt"
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetNoApps(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	w := makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	var apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, 0, len(apps))
}

func TestOneSimpleApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name:       "test",
		AutoDeploy: true,
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Make POST request
	w := makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))
	data := fromJson[map[string]uint](w.Body)
	appView.ID = data["id"]

	// GET request to check if it was inserted
	w = makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	// Check app content + returned id
	var apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, []comm.AppView{appView}, apps)
}

func TestOneComplexApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name:           "test",
		Desc:           "test description",
		SourceURL:      "https://github.com/a/a",
		SourceUsername: "user",
		SourcePassword: "pass",
		AutoDeploy:     true,
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Make POST request
	w := makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))
	data := fromJson[map[string]uint](w.Body)
	appView.ID = data["id"]

	// GET request to check if it was inserted
	w = makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	// Check app content + returned id
	var apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, []comm.AppView{appView}, apps)
}

func TestMultipleApps(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	createRequests := []comm.CreateAppRequest{
		{
			Name:           "test1",
			Desc:           "test1 description",
			SourceURL:      "https://github.com/a/a",
			SourceUsername: "user1",
			SourcePassword: "pass1",
			AutoDeploy:     true,
		},
		// Missing DB fields,
		{
			Name:           "test2",
			SourceUsername: "user2",
			SourcePassword: "pass2",
			AutoDeploy:     false,
		},
		// Missing auth fields
		{
			Name:       "test3",
			Desc:       "test3 description",
			SourceURL:  "https://github.com/c/c",
			AutoDeploy: false,
		},
	}

	var appViews = make([]comm.AppView, len(createRequests))
	for i, createRequest := range createRequests {
		utils.CopyFields(&createRequest, &appViews[i])
	}

	// Make POST requests
	for i, createRequest := range createRequests {
		w := makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(createRequest))
		data := fromJson[map[string]uint](w.Body)
		appViews[i].ID = data["id"]
	}

	// GET request to check if everything was inserted
	w := makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	// Check app content + returned id
	var apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, appViews, apps)
}

func TestRecreateApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name: "test",
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Insert app
	makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))

	// Insert app again
	w := makeRequest(webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))
	assertStatusCode(t, w, 409)
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

	// GET request to check if it was deleted
	w = makeOKRequest(t, webServer, "GET", "/api/v1/applications", nil)

	// Check app content + returned id
	apps = fromJson[[]comm.AppView](w.Body)
	assert.Equal(t, 0, len(apps))
}

func TestDeleteNonexistentApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	w := makeRequest(webServer, "DELETE", "/api/v1/applications/1", nil)
	assertStatusCode(t, w, 404)
}

func TestGetAppByID(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name: "test",
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Insert app
	w := makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))
	data := fromJson[map[string]uint](w.Body)
	appView.ID = data["id"]

	// GET request to check if it was inserted
	w = makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v", appView.ID), nil)

	// Check app content + returned id
	var app = fromJson[comm.AppView](w.Body)
	assert.Equal(t, appView, app)
}

func TestGetAppByName(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name: "test",
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Insert app
	w := makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))
	data := fromJson[map[string]uint](w.Body)
	appView.ID = data["id"]

	// GET request to check if it was inserted
	w = makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v", appView.Name), nil)

	// Check app content + returned id
	var app = fromJson[comm.AppView](w.Body)
	assert.Equal(t, appView, app)
}

func TestGetNonExistingApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	w := makeRequest(webServer, "GET", "/api/v1/applications/1", nil)
	assertStatusCode(t, w, 404)
}

func TestUpdateApp(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	appCreateQuest := comm.CreateAppRequest{
		Name: "test",
	}
	var appView comm.AppView
	utils.CopyFields(&appCreateQuest, &appView)

	// Insert app
	w := makeOKRequest(t, webServer, "POST", "/api/v1/applications", toJson(appCreateQuest))
	data := fromJson[map[string]uint](w.Body)
	appView.ID = data["id"]

	// Update app
	appUpdateQuest := comm.CreateAppRequest{
		Desc:           "test description",
		SourceURL:      "https://github.com/a/a",
		SourceUsername: "user",
		SourcePassword: "pass",
		AutoDeploy:     true,
	}

	// Make PATCH request
	makeOKRequest(t, webServer, "PATCH", fmt.Sprintf("/api/v1/applications/%v", appView.ID), toJson(appUpdateQuest))

	// GET request to check if it was updated
	w = makeOKRequest(t, webServer, "GET", fmt.Sprintf("/api/v1/applications/%v", appView.ID), nil)

	// Check app content + returned id
	var finalApp comm.AppView
	utils.CopyFields(&appUpdateQuest, &finalApp)
	finalApp.ID = appView.ID
	finalApp.Name = "test"

	var app = fromJson[comm.AppView](w.Body)
	assert.Equal(t, finalApp, app)
}
