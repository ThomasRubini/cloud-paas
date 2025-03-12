// Contains the Config singleton that holds the configuration
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	VERBOSE                bool   `env:"VERBOSE" envDefault:"false"`
	REPO_FETCH_ENABLE      bool   `env:"REPO_FETCH_ENABLE" envDefault:"true"`
	REPO_FETCH_PERIOD_SECS int    `env:"REPO_FETCH_PERIOD_SECS"`
	REPO_DIR               string `env:"REPO_DIR"`
	DB_URL                 string `env:"DB_URL"`
	OIDC_BASE_URL          string `env:"OIDC_BASE_URL"`
	OIDC_CLIENT_ID         string `env:"OIDC_CLIENT_ID"`
	OIDC_CLIENT_SECRET     string `env:"OIDC_CLIENT_SECRET"`
	OIDC_USER_ID           string `env:"OIDC_USER_ID"`
	OIDC_USER_PASSWORD     string `env:"OIDC_USER_PASSWORD"`
	OIDC_REALM             string `env:"OIDC_REALM"`

	SECRETS_IMPL      string `env:"SECRETS_IMPL"`
	SECRETS_IMPL_FILE string `env:"SECRETS_IMPL_FILE" envDefault:""`
}

var configInst *Config

// get the application configuration
func Get() Config {
	if configInst == nil {
		panic("config not initialized")
	}
	return *configInst
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		// logrus not setup yet
		fmt.Printf("Did not load .env file (%v)\n", err)
	}

	c := &Config{}
	err = env.ParseWithOptions(c, env.Options{
		RequiredIfNoDef: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Error parsing config: %v", err))
	}
	configInst = c
}
