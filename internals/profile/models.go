package profile

const (
	maxResponseSize = 1 << 20 // 1MB
)

type GenderResponse struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Count       int     `json:"count"`
}

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   *int   `json:"age,omitempty"`
}

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type NationalityResponse struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}
