// Contains the Config singleton that holds the configuration
package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	VERBOSE                bool
	REPO_FETCH_PERIOD_SECS int
	REPO_DIR               string
	DB_URL                 string
	OIDC_BASE_URL          string
	OIDC_CLIENT_ID         string
	OIDC_CLIENT_SECRET     string
	OIDC_USER_ID           string
	OIDC_USER_PASSWORD     string
	OIDC_REALM             string
}

var configInst *Config

// get the application configuration
func Get() Config {
	if configInst == nil {
		panic("config not initialized")
	}
	return *configInst
}

func getEnvInt(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(fmt.Sprintf("Invalid value for %s: %v", key, err))
	}
	return val
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		// logrus not setup yet
		fmt.Printf("Did not load .env file (%v)\n", err)
	}
	configInst = &Config{
		VERBOSE:                !slices.Contains([]string{"false", "0", ""}, os.Getenv("VERBOSE")),
		REPO_FETCH_PERIOD_SECS: getEnvInt("REPO_FETCH_PERIOD_SECS"),
		REPO_DIR:               os.Getenv("REPO_DIR"),
		DB_URL:                 os.Getenv("DB_URL"),
		OIDC_BASE_URL:          os.Getenv("OIDC_BASE_URL"),
		OIDC_USER_ID:           os.Getenv("OIDC_USER_ID"),
		OIDC_USER_PASSWORD:     os.Getenv("OIDC_USER_PASSWORD"),
		OIDC_REALM:             os.Getenv("OIDC_REALM"),
		OIDC_CLIENT_ID:         os.Getenv("OIDC_CLIENT_ID"),
		OIDC_CLIENT_SECRET:     os.Getenv("OIDC_CLIENT_SECRET"),
	}
}
