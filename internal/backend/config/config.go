package config

import (
	"os"
	"slices"
)

type Config struct {
	VERBOSE bool

	DB_URL     string
}

var configInst *Config

func Get() Config {
	if configInst == nil {
		panic("config not initialized")
	}
	return *configInst
}

func Init() {
	configInst = &Config{
		VERBOSE:     !slices.Contains([]string{"false", "0", ""}, os.Getenv("VERBOSE")),
		DB_URL:     os.Getenv("DB_URL"),
	}
}
