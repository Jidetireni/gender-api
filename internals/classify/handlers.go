package classify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("(%d) %s", e.Status, e.Message)
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func encodeError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		encode(w, apiErr.Status, ErrorResponse{
			Status:  "error",
			Message: apiErr.Message,
		})
		return
	}

	// Raw / unexpected error — never leak internals to the caller.
	encode(w, http.StatusInternalServerError, ErrorResponse{
		Status:  "error",
		Message: "internal server error",
	})
}

func HandleClassify(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			encodeError(w, &APIError{Status: http.StatusBadRequest, Message: "name query parameter is required"})
			return
		}

		for _, char := range name {
			if char >= '0' && char <= '9' {
				encodeError(w, &APIError{Status: http.StatusUnprocessableEntity, Message: "name must not contain numbers"})
				return
			}
		}

		result, err := svc.Classify(r.Context(), name)
		if err != nil {
			encodeError(w, err)
			return
		}

		encode(w, http.StatusOK, result)
	}
}
