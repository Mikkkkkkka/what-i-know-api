package config

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func LoadDotEnv(paths ...string) error {
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if err := loadDotEnvFile(path); err != nil {
			return err
		}
	}

	return nil
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		value = strings.TrimSpace(value)
		value = trimInlineComment(value)
		value = parseEnvValue(value)

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func trimInlineComment(value string) string {
	if len(value) == 0 {
		return value
	}
	if value[0] == '"' || value[0] == '\'' {
		return value
	}

	commentIndex := strings.Index(value, " #")
	if commentIndex >= 0 {
		return strings.TrimSpace(value[:commentIndex])
	}

	return value
}

func parseEnvValue(value string) string {
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			unquoted, err := strconv.Unquote(value)
			if err == nil {
				return unquoted
			}

			return value[1 : len(value)-1]
		}
	}

	return value
}
