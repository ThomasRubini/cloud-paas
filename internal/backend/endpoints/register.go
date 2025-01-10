package endpoints

import (
	"bytes"
	"cloud-paas/internal/backend/config"
	"cloud-paas/internal/noerror"
	"cloud-paas/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TODO maybe use raw Keycloak API rather than gocloak just for that ?
func getKeycloakAccessToken() (string, error) {
	cfg := config.Get()
	client := gocloak.NewClient(cfg.OIDC_BASE_URL)

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
	statusCode               int
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

	url, err := url.JoinPath(cfg.OIDC_BASE_URL, fmt.Sprintf("/admin/realms/%v/users", cfg.OIDC_REALM))
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

	if utils.IsStatusCodeOk(resp.StatusCode) {
		return registerRequestResult{statusCode: resp.StatusCode}
	} else if resp.StatusCode == http.StatusBadRequest {
		return registerRequestResult{statusCode: resp.StatusCode, keycloackValidationError: noerror.ReadAll(resp.Body)}
	} else {
		return registerRequestResult{statusCode: resp.StatusCode, err: fmt.Errorf("got error response. Body: %v", noerror.ReadAll(resp.Body))}
	}
}

func translateKeycloakError(err string) string {
	var keycloakError struct {
		ErrorMessage string        `json:"errorMessage"`
		Field        string        `json:"field"`
		Params       []interface{} `json:"params"`
	}

	if err := json.Unmarshal([]byte(err), &keycloakError); err != nil {
		return "Unknown error"
	}

	switch keycloakError.ErrorMessage {
	case "error-invalid-length":
		return fmt.Sprintf("Attribute %v must be between %v and %v characters long", keycloakError.Field, keycloakError.Params[1], keycloakError.Params[2])
	default:
		logrus.Error("Register failed because of unknown Keycloak error: ", keycloakError)
		return "Unknown error"
	}
}

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register godoc
// @Summary      Register an account
// @Produce      json
// @Success      200 {object} RegisterInput
// @Router       /api/v1/register [post]
func register(c *gin.Context) {
	var req RegisterInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	ret := makeRegisterRequest(req.Username, req.Password)
	if ret.statusCode == 0 {
		ret.statusCode = 500
	}
	if ret.err != nil {
		c.JSON(ret.statusCode, gin.H{"error": ret.err.Error()})
		return
	} else if ret.keycloackValidationError != "" {
		c.JSON(ret.statusCode, gin.H{"status": translateKeycloakError(ret.keycloackValidationError)})
	} else {
		c.JSON(200, gin.H{"status": "user registered successfully"})
	}
}

func initRegister(g *gin.RouterGroup) {
	g.POST("/register", register)
}
