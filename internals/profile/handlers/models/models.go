package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

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

type APIResponse[T any] struct {
	Status  string  `json:"status"`
	Message *string `json:"message,omitempty"`
	Count   *int    `json:"count,omitempty"`
	Page    *int    `json:"page,omitempty"`
	Limit   *int    `json:"limit,omitempty"`
	Total   *int    `json:"total,omitempty"`
	Data    T       `json:"data"`
}

type CreateProfileRequest struct {
	Name string `json:"name"`
}

type Profile struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Gender             string    `json:"gender"`
	GenderProbability  float64   `json:"gender_probability"`
	Age                int       `json:"age"`
	AgeGroup           string    `json:"age_group"`
	CountryID          string    `json:"country_id"`
	CountryName        string    `json:"country_name"`
	CountryProbability float64   `json:"country_probability"`
	CreatedAt          time.Time `json:"created_at"`
}

type ProfileShort struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Age       int       `json:"age"`
	AgeGroup  string    `json:"age_group"`
	CountryID string    `json:"country_id"`
}

type ListResult[T any] struct {
	Total int64 `json:"total,omitempty"`
	Data  []*T  `json:"data"`
}
