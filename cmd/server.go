package main

import (
	"net/http"

	"github.com/Jidetireni/gender-api/config"
	"github.com/Jidetireni/gender-api/internals/pkg/database/postgres"
	"github.com/Jidetireni/gender-api/internals/profile"
	"github.com/Jidetireni/gender-api/internals/profile/handlers"
	"github.com/go-chi/chi/v5"
)

func NewServer(
	cfg *config.Config,
	profileSvc *profile.Service,
	postgresDB *postgres.PostgresDB,
) http.Handler {
	r := chi.NewRouter()
	addRoutes(
		r,
		profileSvc,
	)

	return r
}

func addRoutes(
	r *chi.Mux,
	profileSvc *profile.Service,
) {

	r.Route("/api", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				next.ServeHTTP(w, req)
			})
		})
		r.Post("/profiles", handlers.HandleCreateProfile(profileSvc))
		r.Get("/profiles/{id}", handlers.HandleGetProfile(profileSvc))
		r.Get("/profiles", handlers.HandleListProfiles(profileSvc))
		r.Delete("/profiles/{id}", handlers.HandleDeleteProfile(profileSvc))
	})

}
