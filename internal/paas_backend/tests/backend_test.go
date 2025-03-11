package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
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

func TestHealth(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState())

	// Test health URL
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	webServer.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
}
