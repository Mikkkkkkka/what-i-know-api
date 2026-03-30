package config

import (
	"os"
	"strconv"
	"time"
)

type HTTPConfig struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	APIBasePath     string
}

func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Address:         ":8080",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     30 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		APIBasePath:     "/",
	}
}

func LoadHTTPConfigFromEnv() HTTPConfig {
	return loadHTTPConfigFromEnv()
}

func loadHTTPConfigFromEnv() HTTPConfig {
	cfg := DefaultHTTPConfig()

	if value := os.Getenv("HTTP_ADDRESS"); value != "" {
		cfg.Address = value
	}
	if value := os.Getenv("HTTP_API_BASE_PATH"); value != "" {
		cfg.APIBasePath = value
	}
	if value := os.Getenv("HTTP_READ_TIMEOUT_SECONDS"); value != "" {
		cfg.ReadTimeout = durationSecondsOrDefault(value, cfg.ReadTimeout)
	}
	if value := os.Getenv("HTTP_WRITE_TIMEOUT_SECONDS"); value != "" {
		cfg.WriteTimeout = durationSecondsOrDefault(value, cfg.WriteTimeout)
	}
	if value := os.Getenv("HTTP_IDLE_TIMEOUT_SECONDS"); value != "" {
		cfg.IdleTimeout = durationSecondsOrDefault(value, cfg.IdleTimeout)
	}
	if value := os.Getenv("HTTP_SHUTDOWN_TIMEOUT_SECONDS"); value != "" {
		cfg.ShutdownTimeout = durationSecondsOrDefault(value, cfg.ShutdownTimeout)
	}

	return cfg
}

func durationSecondsOrDefault(raw string, fallback time.Duration) time.Duration {
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return fallback
	}

	return time.Duration(seconds) * time.Second
}
