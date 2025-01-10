package cli

import (
	"bytes"
	"cloud-paas/internal/cli/config"
	"cloud-paas/internal/noerror"
	"cloud-paas/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"syscall"

	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// Does not parse response
// Has 3 return possibilities:
// - Technical error (returns "", err)
// - Validation failure in backend (returns "validation error", nil)
// - Success (returns "", nil)
func makeRegisterAccountRequest(user, password string) (string, error) {
	url, err := url.JoinPath(config.Get().BACKEND_URL, "/api/v1/register")
	if err != nil {
		return "", fmt.Errorf("failed to join url: %w", err)
	}

	data := map[string]any{
		"username": user,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to parse into json: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	if utils.IsStatusCodeOk(resp.StatusCode) {
		// Success
		return "", nil
	} else if resp.StatusCode == http.StatusBadRequest {
		// Parse validation error msg
		type validationError struct {
			Status string `json:"status"`
		}

		var vErr validationError
		if err := json.NewDecoder(resp.Body).Decode(&vErr); err != nil {
			return "", fmt.Errorf("failed to decode validation error response: %w", err)
		}
		return vErr.Status, nil
	} else if resp.StatusCode == http.StatusConflict {
		return "username already exists", nil
	} else {
		return "", fmt.Errorf("unexpected status code: %v. Backend response: %v", resp.StatusCode, noerror.ReadAll(resp.Body))
	}
}

func getUserAndPassword(c *cli.Command) (user string, password string, err error) {
	user = c.String("user")
	if user == "" {
		fmt.Print("Enter username: ")
		_, err = fmt.Scanln(&user)
		if err != nil {
			err = fmt.Errorf("failed to read username: %w", err)
			return
		}
	}

	password = c.String("password")
	if password == "" {
		fmt.Print("Enter password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", fmt.Errorf("failed to read password: %w", err)
		}
		password = string(bytePassword)
		fmt.Println()
	}

	return user, password, nil
}

// Make login request against keycloak
func makeLoginRequest(user, password string) (string, error) {
	cfg := config.Get()

	res, err := url.JoinPath(cfg.OIDC_REALM_URL, "/protocol/openid-connect/token")
	if err != nil {
		return "", fmt.Errorf("failed to join url: %w", err)
	}

	resp, err := http.PostForm(res, url.Values{
		"grant_type": {"password"},
		"client_id":  {"paas-cli"},
		"username":   {user},
		"password":   {password},
	})
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d (%v)", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return token, nil
}

func RegisterAction(ctx context.Context, c *cli.Command) error {
	conf := config.Get()
	if conf.AUTH_TOKEN != "" {
		return fmt.Errorf("already logged in")
	}

	user, password, err := getUserAndPassword(c)
	if err != nil {
		return fmt.Errorf("failed to read user and password: %w", err)
	}

	validationError, err := makeRegisterAccountRequest(user, password)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	if validationError != "" {
		fmt.Printf("Validation error: %v\n", validationError)
		return nil
	}

	token, err := makeLoginRequest(user, password)
	if err != nil {
		return fmt.Errorf("registration succeeded, but login to get token failed: %w", err)
	}

	fmt.Printf("Access token: %v\n", token)

	panic("TODO: save config")
}
