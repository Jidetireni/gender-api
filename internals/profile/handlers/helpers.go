package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Jidetireni/gender-api/internals/profile/handlers/models"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z]+$`)

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func encodeError(w http.ResponseWriter, err error) {
	var apiErr *models.APIError
	if errors.As(err, &apiErr) {
		statusStr := "error"
		if apiErr.Status == http.StatusBadGateway {
			statusStr = "502"
		}
		encode(w, apiErr.Status, models.ErrorResponse{
			Status:  statusStr,
			Message: apiErr.Message,
		})
		return
	}

	log.Printf("unexpected error: %v", err.Error())
	// Raw / unexpected error — never leak internals to the caller.
	encode(w, http.StatusInternalServerError, models.ErrorResponse{
		Status:  "error",
		Message: "internal server error",
	})
	log.Printf("unexpected error: %v", err)
}
