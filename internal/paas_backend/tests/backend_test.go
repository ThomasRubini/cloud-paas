package tests

import (
	"testing"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	webServer := paas_backend.SetupWebServer(fakeState(t))

	// Test health URL
	w := makeOKRequest(t, webServer, "GET", "/health", nil)
	assert.Equal(t, "OK", toString(w.Body))

	w = makeRequest(webServer, "GET", "/nonexistent", nil)
	assert.Equal(t, 404, w.Code)
}
