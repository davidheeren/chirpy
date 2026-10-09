package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if (authHeader == "") {
		return "", errors.New("Authorization header cannot be empty")
	}

	// remove Bearer prefix + extra whitespace
	jwt := strings.TrimLeft(authHeader, "Bearer")
	if (jwt == authHeader) {
		return "", errors.New("'Bearer' must start the Authorization header")
	}
	jwt = strings.TrimSpace(jwt)

	return jwt, nil
}
