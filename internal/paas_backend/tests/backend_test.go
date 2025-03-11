package tests

import (
	"io"
	"net/http/httptest"
	"testing"

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

	return utils.State{
		Db: gorm_db,
	}
}

func makeRequest(webServer *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)
	webServer.ServeHTTP(w, r)
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
