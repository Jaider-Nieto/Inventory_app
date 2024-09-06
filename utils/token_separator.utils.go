package utils

import (
	"errors"
	"strings"
)

func TokenSeparator(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization format")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
