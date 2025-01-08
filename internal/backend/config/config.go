package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

type Config struct {
	VERBOSE bool

	DB_URL string
}

var configInst *Config

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
	configInst = &Config{
		VERBOSE: !slices.Contains([]string{"false", "0", ""}, os.Getenv("VERBOSE")),
		DB_URL:  os.Getenv("DB_URL"),
	}
}
