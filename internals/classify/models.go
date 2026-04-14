package classify

import "time"

const (
	maxResponseSize = 1 << 20 // 1MB
)

type GenderResponse struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type APIResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type ProcessedData struct {
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Probability float64   `json:"probability"`
	SampleSize  int       `json:"sample_size"`
	IsConfident bool      `json:"is_confident"`
	ProcessedAt time.Time `json:"processed_at"`
}
