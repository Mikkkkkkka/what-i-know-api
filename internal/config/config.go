package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseDSN         string
	HTTPAddress         string
	HTTPAPIBasePath     string
	HTTPReadTimeout     time.Duration
	HTTPWriteTimeout    time.Duration
	HTTPIdleTimeout     time.Duration
	HTTPShutdownTimeout time.Duration
}

func Load() Config {
	cfg := Config{
		DatabaseDSN:         os.Getenv("DATABASE_DSN"),
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

func durationSecondsOrDefault(raw string, fallback time.Duration) time.Duration {
	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		return fallback
	}

	return time.Duration(seconds) * time.Second
}
