package classify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Jidetireni/gender-api/config"
)

type Service struct {
	httpClient           *http.Client
	genderizedAPIBaseURL string
}

func New(config *config.Config) *Service {
	return &Service{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		genderizedAPIBaseURL: config.GenderizedAPIBaseURL,
	}
}

func (s *Service) getGenderResponse(ctx context.Context, name string) (*GenderResponse, error) {
	url := fmt.Sprintf("%s/?name=%s", s.genderizedAPIBaseURL, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, err
	}

	genderResponse := &GenderResponse{}
	if err := json.Unmarshal(respBody, genderResponse); err != nil {
		return nil, err
	}

	return genderResponse, nil
}

func (s *Service) Classify(ctx context.Context, name string) (*APIResponse[ProcessedData], error) {
	genderResponse, err := s.getGenderResponse(ctx, name)
	if err != nil {
		return nil, err
	}

	if genderResponse == nil || genderResponse.Count == 0 || genderResponse.Gender == "" {
		return nil, &APIError{
			Status:  http.StatusNotFound,
			Message: "No prediction available for the provided name",
		}
	}

	sampleSize := genderResponse.Count
	isConfident := genderResponse.Probability >= 0.7 && sampleSize >= 100

	return &APIResponse[ProcessedData]{
		Status: "success",
		Data: ProcessedData{
			Name:        name,
			Gender:      genderResponse.Gender,
			Probability: genderResponse.Probability,
			SampleSize:  sampleSize,
			IsConfident: isConfident,
			ProcessedAt: time.Now().UTC(),
		},
	}, nil

}
