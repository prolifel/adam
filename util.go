package main

import (
	"errors"
	"net/http"
	"strings"
)

func validateBearer(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("unauthorized")
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 {
		return "", errors.New("unauthorized")
	}

	token := strings.TrimSpace(parts[1])
	if len(token) == 0 {
		return "", errors.New("unauthorized")
	}

	return token, nil
}

func validateToken(r *http.Request, token string) error {
	bearer, err := validateBearer(r)
	if err != nil {
		return err
	}

	if token != bearer {
		return errors.New("unauthorized")
	}

	return nil
}
