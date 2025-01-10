package cli

import (
	"bytes"
	"cloud-paas/internal/cli/config"
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

func registerAccountRequest(user, password string) (*string, error) {
	url, err := url.JoinPath(config.Get().BackendURL, "/api/v1/register")
	if err != nil {
		return nil, fmt.Errorf("failed to join url: %w", err)
	}

	data := map[string]any{
		"username": user,
		"password": password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse into json: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d (%v)", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return nil, fmt.Errorf("access token not found in response")
	}

	fmt.Println("Access Token:", token)
	return &token, nil
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

	token, err := registerAccountRequest(user, password)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	conf.AuthToken = *token
	panic("TODO: save config")
}
