package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	HTTPAddress         string
	HTTPAPIBasePath     string
	HTTPReadTimeout     time.Duration
	HTTPWriteTimeout    time.Duration
	HTTPIdleTimeout     time.Duration
	HTTPShutdownTimeout time.Duration
}

func Load() Config {
	cfg := Config{
		DBHost:              os.Getenv("DBHost"),
		DBPort:              os.Getenv("DBPort"),
		DBUser:              os.Getenv("DBUser"),
		DBPassword:          os.Getenv("DBPassword"),
		DBName:              os.Getenv("DBName"),
		DBSSLMode:           os.Getenv("DBSSLMode"),
		HTTPAddress:         os.Getenv("HTTP_ADDRESS"),
		HTTPAPIBasePath:     os.Getenv("HTTP_API_BASE_PATH"),
		HTTPReadTimeout:     5 * time.Second,
		HTTPWriteTimeout:    10 * time.Second,
		HTTPIdleTimeout:     30 * time.Second,
		HTTPShutdownTimeout: 10 * time.Second,
	}

	if cfg.HTTPAddress == "" {
		cfg.HTTPAddress = ":8080"
	}
	if cfg.DBPort == "" {
		cfg.DBPort = "5432"
	}
	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "disable"
	}
	if cfg.HTTPAPIBasePath == "" {
		cfg.HTTPAPIBasePath = "/"
	}
	if value := os.Getenv("HTTP_READ_TIMEOUT_SECONDS"); value != "" {
		cfg.HTTPReadTimeout = durationSecondsOrDefault(value, cfg.HTTPReadTimeout)
	}
	if value := os.Getenv("HTTP_WRITE_TIMEOUT_SECONDS"); value != "" {
		cfg.HTTPWriteTimeout = durationSecondsOrDefault(value, cfg.HTTPWriteTimeout)
	}
	if value := os.Getenv("HTTP_IDLE_TIMEOUT_SECONDS"); value != "" {
		cfg.HTTPIdleTimeout = durationSecondsOrDefault(value, cfg.HTTPIdleTimeout)
	}
	if value := os.Getenv("HTTP_SHUTDOWN_TIMEOUT_SECONDS"); value != "" {
		cfg.HTTPShutdownTimeout = durationSecondsOrDefault(value, cfg.HTTPShutdownTimeout)
	}

	return cfg
}

func (c Config) MissingRequiredDBEnv() []string {
	var missing []string

	if c.DBHost == "" {
		missing = append(missing, "DBHost")
	}
	if c.DBUser == "" {
		missing = append(missing, "DBUser")
	}
	if c.DBPassword == "" {
		missing = append(missing, "DBPassword")
	}
	if c.DBName == "" {
		missing = append(missing, "DBName")
	}

	return missing
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)
}

func FormatMissingEnv(names []string) string {
	return strings.Join(names, ", ")
}

func durationSecondsOrDefault(raw string, fallback time.Duration) time.Duration {
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return fallback
	}

	return time.Duration(seconds) * time.Second
}
