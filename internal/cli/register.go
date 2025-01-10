package cli

import (
	"bytes"
	"cloud-paas/internal/cli/config"
	"cloud-paas/internal/noerror"
	"context"
	"encoding/json"
	"fmt"
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
	url, err := url.JoinPath(config.Get().BackendURL, "/api/v1/register")
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

	if resp.StatusCode == http.StatusOK {
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
	} else {
		return "", fmt.Errorf("unexpected status code: %v (%v)", resp.StatusCode, noerror.ReadAll(resp.Body))
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

func RegisterAction(ctx context.Context, c *cli.Command) error {
	conf := config.Get()
	if conf.AuthToken != "" {
		return fmt.Errorf("already logged in")
	}

	user, password, err := getUserAndPassword(c)
	if err != nil {
		return fmt.Errorf("failed to read user and password: %w", err)
	}

	fmt.Printf("%v %v\n", user, password)

	validationError, err := makeRegisterAccountRequest(user, password)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	if validationError != "" {
		fmt.Printf("Validation error: %v\n", validationError)
		return nil
	}

	panic("TODO: save config")
}
