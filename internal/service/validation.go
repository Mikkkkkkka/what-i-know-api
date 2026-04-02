package service

import (
	"errors"
	"strings"
)

var ErrInvalidInput = errors.New("invalid input")

func normalizeRequiredString(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return "", ErrInvalidInput
	}

	return normalized, nil
}
