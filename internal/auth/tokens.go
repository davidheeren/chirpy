package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header cannot be empty")
	}

	// remove Bearer prefix + extra whitespace
	jwt := strings.TrimLeft(authHeader, "Bearer")
	if jwt == authHeader {
		return "", errors.New("'Bearer' must start the Authorization header")
	}
	jwt = strings.TrimSpace(jwt)

	return jwt, nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	// Note that no error handling is necessary, as Read always succeeds.
	rand.Read(token)
	return hex.EncodeToString(token)
}
