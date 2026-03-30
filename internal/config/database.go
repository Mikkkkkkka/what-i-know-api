package config

import (
	"os"
	"strings"
)

type DatabaseConfig struct {
	DSN string
}

func LoadDatabaseConfigFromEnv() DatabaseConfig {
	return loadDatabaseConfigFromEnv()
}

func loadDatabaseConfigFromEnv() DatabaseConfig {
	return DatabaseConfig{
		DSN: strings.TrimSpace(os.Getenv("DATABASE_DSN")),
	}
}
