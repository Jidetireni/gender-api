package config

import (
	"log"
	"os"
)

type Config struct {
	Host                 string
	Port                 string
	Env                  string
	GenderizedAPIBaseURL string
}

func validateENV() {
	environmentVariables := []string{
		"HOST",
		"PORT",
		"ENV",
		"GENDERIZED_API_BASE_URL",
	}

	for _, env := range environmentVariables {
		if os.Getenv(env) == "" {
			log.Fatalf("Environment variable %s is not set", env)
		}
	}
}

func New() *Config {
	validateENV()

	return &Config{
		Host:                 os.Getenv("HOST"),
		Port:                 os.Getenv("PORT"),
		Env:                  os.Getenv("ENV"),
		GenderizedAPIBaseURL: os.Getenv("GENDERIZED_API_BASE_URL"),
	}
}
