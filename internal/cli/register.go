package cli

import (
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

func makeDirectGrantRequest(issuer_url, user, password string) (*string, error) {
	res, err := url.JoinPath(issuer_url, "/protocol/openid-connect/token")
	if err != nil {
		return nil, fmt.Errorf("failed to join url: %w", err)
	}

	resp, err := http.PostForm(res, url.Values{
		"grant_type": {"password"},
		"client_id":  {"cli"},
		"username":   {user},
		"password":   {password},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

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

	token, err := makeDirectGrantRequest(conf.OIDCIssuerURL, user, password)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	conf.AuthToken = *token
	panic("TODO: save config")
}
