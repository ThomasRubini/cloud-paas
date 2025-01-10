package endpoints

import (
	"bytes"
	"cloud-paas/internal/backend/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

// TODO maybe use raw Keycloak API rather than gocloak just for that ?
func getKeycloakAccessToken() (string, error) {
	cfg := config.Get()
	client := gocloak.NewClient(cfg.OIDC_URL)

	token, err := client.Login(
		context.Background(),
		cfg.OIDC_CLIENT_ID, cfg.OIDC_CLIENT_SECRET,
		cfg.OIDC_REALM,
		cfg.OIDC_USER_ID, cfg.OIDC_USER_PASSWORD,
	)
	if err != nil {
		return "", fmt.Errorf("failed to login to keycloak: %w", err)
	}

	return token.AccessToken, nil
}

type registerRequestResult struct {
	err                      error
	keycloackValidationError string
}

// We make the request without gocloak because it doesn't give us all the response, which we need
func makeRegisterRequest(username, password string) registerRequestResult {
	cfg := config.Get()
	accessToken, err := getKeycloakAccessToken()
	if err != nil {
		return registerRequestResult{err: fmt.Errorf("failed to get Keycloak access token: %w", err)}
	}

	url, err := url.JoinPath(cfg.OIDC_URL, fmt.Sprintf("/admin/realms/%v/users", cfg.OIDC_REALM))
	if err != nil {
		return registerRequestResult{err: fmt.Errorf("failed to create Keycloak request URL: %w", err)}
	}

	user := gocloak.User{
		Username: &username,
		Enabled:  gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     &password,
			},
		},
	}

	b, err := json.Marshal(user)
	if err != nil {
		return registerRequestResult{err: fmt.Errorf("failed to convert Keycloak UserRepresentation to json: %w", err)}
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return registerRequestResult{err: fmt.Errorf("failed to create Keycloak request: %w", err)}
	}

	req.Header.Add("authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return registerRequestResult{err: fmt.Errorf("failed to send Keycloak request: %w", err)}
	}

	if resp.StatusCode == http.StatusOK {
		return registerRequestResult{}
	} else if resp.StatusCode == http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return registerRequestResult{keycloackValidationError: string(b)}
	} else {
		b, _ := io.ReadAll(resp.Body)
		return registerRequestResult{err: fmt.Errorf("unexpected status code: %d (%v)", resp.StatusCode, string(b))}
	}
}

func register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	ret := makeRegisterRequest(req.Username, req.Password)
	if ret.err != nil {
		c.JSON(500, gin.H{"error": ret.err.Error()})
		return
	} else if ret.keycloackValidationError != "" {
		c.Data(400, "application/json", []byte(ret.keycloackValidationError))
	} else {
		c.JSON(200, gin.H{"status": "user registered successfully"})
	}
}

func initRegister(g *gin.RouterGroup) {
	g.POST("/register", register)
}
