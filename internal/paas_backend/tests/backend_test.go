package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/comm"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func fakeState() utils.State {
	gorm_db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		panic(err)
	}

	if err := paas_backend.MigrateModels(gorm_db); err != nil {
		panic(err)
	}

	return utils.State{
		Db: gorm_db,
	}
}

func makeRequest(webServer *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, body)
	webServer.ServeHTTP(w, r)
	return w
}

func makeOKRequest(t *testing.T, webServer *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	w := makeRequest(webServer, method, path, body)
	if w.Code != 200 {
		fmt.Println("Request failed !")
		fmt.Println("Body:")
		fmt.Println(w.Body)
		assert.Equal(t, 200, w.Code)
	}
	return w
}

func TestHealth(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	// Test health URL
	w := makeRequest(webServer, "GET", "/health", nil)
	assert.Equal(t, 200, w.Code)

	w = makeRequest(webServer, "GET", "/nonexistent", nil)
	assert.Equal(t, 404, w.Code)
}

func toJson(v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}

func fromJson(body io.Reader, v interface{}) {
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		panic(err)
	}
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
