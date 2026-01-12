package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	databaseURL, err := getEnv("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}

	port, err := getEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return Config{
		DatabaseURL: databaseURL,
		ServerPort:  port,
	}, nil
}

func getEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return val, nil
}
