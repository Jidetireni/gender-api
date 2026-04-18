package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                  string
	Port                  string
	Env                   string
	GenderizedAPIBaseURL  string
	AgifyAPIBaseURL       string
	NationalizeAPIBaseURL string
	DBURL                 string
	RedisURL              string
}

func validateENV() {
	environmentVariables := []string{
		"HOST",
		"PORT",
		"ENV",
		"GENDERIZED_API_BASE_URL",
		"AGIFY_API_BASE_URL",
		"NATIONAIZE_API_BASE_URL",
		"DB_URL",
		"REDIS_URL",
	}

	_ = godotenv.Load()

	for _, env := range environmentVariables {
		if os.Getenv(env) == "" {
			log.Fatalf("Environment variable %s is not set", env)
		}
	}
}

func New() *Config {
	validateENV()

	return &Config{
		Host:                  os.Getenv("HOST"),
		Port:                  os.Getenv("PORT"),
		Env:                   os.Getenv("ENV"),
		GenderizedAPIBaseURL:  os.Getenv("GENDERIZED_API_BASE_URL"),
		AgifyAPIBaseURL:       os.Getenv("AGIFY_API_BASE_URL"),
		NationalizeAPIBaseURL: os.Getenv("NATIONAIZE_API_BASE_URL"),
		DBURL:                 os.Getenv("DB_URL"),
		RedisURL:              os.Getenv("REDIS_URL"),
	}
}
