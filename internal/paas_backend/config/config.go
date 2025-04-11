// Contains the Config singleton that holds the configuration
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	VERBOSE                string `env:"VERBOSE" envDefault:""`
	REPO_FETCH_ENABLE      bool   `env:"REPO_FETCH_ENABLE" envDefault:"true"`
	REPO_FETCH_PERIOD_SECS int    `env:"REPO_FETCH_PERIOD_SECS"`
	REPO_DIR               string `env:"REPO_DIR"`
	DB_URL                 string `env:"DB_URL"`
	DB_PASSWORD            string `env:"DB_PASSWORD"`
	OIDC_BASE_URL          string `env:"OIDC_BASE_URL"`
	OIDC_CLIENT_ID         string `env:"OIDC_CLIENT_ID"`
	OIDC_CLIENT_SECRET     string `env:"OIDC_CLIENT_SECRET"`
	OIDC_USER_ID           string `env:"OIDC_USER_ID"`
	OIDC_USER_PASSWORD     string `env:"OIDC_USER_PASSWORD"`
	OIDC_REALM             string `env:"OIDC_REALM"`

	REGISTRY_REPO_URI string `env:"REGISTRY_REPO_URI"`
	REGISTRY_USER     string `env:"REGISTRY_USER" envDefault:""`
	REGISTRY_PASSWORD string `env:"REGISTRY_PASSWORD" envDefault:""`

	SECRETS_IMPL      string `env:"SECRETS_IMPL"`
	SECRETS_IMPL_FILE string `env:"SECRETS_IMPL_FILE" envDefault:""`

	KUBE_NAMESPACE_PREFIX string `env:"KUBE_NAMESPACE_PREFIX"`
	REGISTRY_TAG_PREFIX   string `env:"REGISTRY_TAG_PREFIX"`
}

var configInst *Config

// get the application configuration
func Get() *Config {
	if configInst == nil {
		panic("config not initialized")
	}
	return configInst
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
