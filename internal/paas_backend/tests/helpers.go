package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/config"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/secretsprovider"
	"github.com/ThomasRubini/cloud-paas/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func fakeState(t *testing.T) utils.State {
	gorm_db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		panic(err)
	}

	if err := paas_backend.MigrateModels(gorm_db); err != nil {
		panic(err)
	}

	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		panic(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return utils.State{
		Config:          &config.Config{},
		Db:              gorm_db,
		SecretsProvider: secretsprovider.Helper{Core: secretsprovider.FromFile(filepath.Join(tmpDir, "/secrets.json"))},
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
		fmt.Println("Request failed ! (expected 200, got", w.Code, "). Body:")
		fmt.Println(w.Body)
		assert.Equal(t, 200, w.Code)
	}
	return w
}

func toJson(v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}

func fromJson[T any](body io.Reader) T {
	var v T
	err := json.NewDecoder(body).Decode(&v)
	if err != nil {
		panic(err)
	}
	return v
}

func toString(body io.Reader) string {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		panic(err)
	}
	return buf.String()
}

func assertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	if w.Code != expected {
		fmt.Printf("Request failed ! (expected %v, got %v). Body:", expected, w.Code)
		fmt.Println(w.Body)
		assert.Equal(t, expected, w.Code)
	}
}
