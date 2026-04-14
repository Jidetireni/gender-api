package main

import (
	"net/http"

	"github.com/Jidetireni/gender-api/config"
	"github.com/Jidetireni/gender-api/internals/classify"
	"github.com/go-chi/chi/v5"
)

func NewServer(cfg *config.Config, classifySvc *classify.Service) http.Handler {
	r := chi.NewRouter()
	addRoutes(
		r,
		classifySvc,
	)

	return r
}

func addRoutes(
	r *chi.Mux,
	classifySvc *classify.Service,
) {

	r.Route("/api", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				next.ServeHTTP(w, req)
			})
		})
		r.Get("/classify", classify.HandleClassify(classifySvc))
	})

}
