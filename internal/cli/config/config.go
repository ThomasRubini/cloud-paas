package config

import (
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// List possible config dirs, in order
func listConfigDirs() []string {
	paths := []string{}
	paths = append(paths, "./")
	paths = append(paths, os.Getenv("XDG_CONFIG_HOME"))
	paths = append(paths, "~/.config/")

	return paths
}

// Find the first config file that exists
func findConfigFile() *string {
	for _, folder := range listConfigDirs() {
		filepath := path.Join(folder, "paas_cli_config.yml")
		if _, err := os.Stat(filepath); err == nil {
			return &filepath
		}
	}
	return nil
}

type Config struct {
	BackendURL     string `yaml:"backend_url"`
	OIDC_REALM_URL string `yaml:"oidc_realm_url"`
	AuthToken      string `yaml:"auth_token"`
}

var configInst *Config

func Get() Config {
	if configInst == nil {
		panic("Config not initialized")
	} else {
		return *configInst
	}
}

func Init() {
	configFile := findConfigFile()
	if configFile == nil {
		return
	}

	data, err := os.ReadFile(*configFile)
	if err != nil {
		return
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	configInst = &cfg
}
