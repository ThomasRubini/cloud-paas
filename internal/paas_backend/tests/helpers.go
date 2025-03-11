package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/secretsprovider"
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
		Db:              gorm_db,
		SecretsProvider: secretsprovider.Helper{Core: secretsprovider.FromFile("/tmp/secrets.json")},
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
