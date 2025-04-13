package config

import (
	"fmt"
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
	BACKEND_URL string `yaml:"backend_url"`
}

var configInst *Config

func Get() Config {
	if configInst == nil {
		panic("Config not initialized")
	} else {
		return *configInst
	}
}

func Save(cfg Config) error {

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config into yaml: %w", err)
	}

	configFile := findConfigFile()
	if configFile == nil {
		return fmt.Errorf("config file not found")
	}

	err = os.WriteFile(*configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	configInst = &cfg
	return nil
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
