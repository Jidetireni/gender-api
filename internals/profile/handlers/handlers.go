package handlers

import (
	"net/http"

	"github.com/Jidetireni/gender-api/internals/profile"
	"github.com/Jidetireni/gender-api/internals/profile/handlers/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func HandleCreateProfile(svc *profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := decode[models.CreateProfileRequest](r)
		if err != nil {
			encodeError(w, &models.APIError{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid type",
			})
			return
		}

		if data.Name == "" {
			encodeError(w, &models.APIError{
				Status:  http.StatusBadRequest,
				Message: "Missing or empty name",
			})
			return
		}

		if !nameRegex.MatchString(data.Name) {
			encodeError(w, &models.APIError{
				Status:  http.StatusUnprocessableEntity,
				Message: "name must only contain alphabetic characters",
			})
			return
		}

		profile, isNew, err := svc.Create(r.Context(), data.Name)
		if err != nil {
			if err.Error() == "Genderize returned an invalid response" ||
				err.Error() == "Agify returned an invalid response" ||
				err.Error() == "Nationalize returned an invalid response" {
				encodeError(w, &models.APIError{
					Status:  http.StatusBadGateway,
					Message: err.Error(),
				})
				return
			}
			encodeError(w, err)
			return
		}

		var message *string
		statusCode := http.StatusCreated
		if !isNew {
			msg := "Profile already exists"
			message = &msg
			statusCode = http.StatusOK
		}

		encode(w, statusCode, models.APIResponse[*models.Profile]{
			Status:  "success",
			Message: message,
			Data:    profile,
		})
	}
}

func HandleGetProfile(svc *profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID := chi.URLParam(r, "id")

		if profileID == "" {
			encodeError(w, &models.APIError{
				Status:  http.StatusBadRequest,
				Message: "profile ID is required",
			})
			return
		}

		id, err := uuid.Parse(profileID)
		if err != nil {
			encodeError(w, &models.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid profile ID",
			})
			return
		}

		profile, err := svc.Get(r.Context(), id)
		if err != nil {
			if err.Error() == "Profile not found" {
				encodeError(w, &models.APIError{
					Status:  http.StatusNotFound,
					Message: err.Error(),
				})
				return
			}
			encodeError(w, err)
			return
		}

		encode(w, http.StatusOK, models.APIResponse[*models.Profile]{
			Status: "success",
			Data:   profile,
		})
	}
}

func HandleListProfiles(svc *profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		filter, err := parseProfileFilters(query)
		if err != nil {
			encodeError(w, &models.APIError{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid query parameters",
			})
			return
		}

		queryOptions, err := parseQueryOptions(query)
		if err != nil {
			encodeError(w, &models.APIError{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid query parameters",
			})
			return
		}

		result, err := svc.List(r.Context(), filter, queryOptions)
		if err != nil {
			encodeError(w, err)
			return
		}

		encode(w, http.StatusOK, models.APIResponse[any]{
			Status: "success",
			Page:   lo.ToPtr(int(queryOptions.Page)),
			Limit:  lo.ToPtr(int(queryOptions.Limit)),
			Total:  lo.ToPtr(int(result.Total)),
			Data:   result.Data,
		})
	}
}

func HandleDeleteProfile(svc *profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID := chi.URLParam(r, "id")

		if profileID == "" {
			encodeError(w, &models.APIError{
				Status:  http.StatusBadRequest,
				Message: "profile ID is required",
			})
			return
		}

		id, err := uuid.Parse(profileID)
		if err != nil {
			encodeError(w, &models.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid profile ID",
			})
			return
		}

		err = svc.Delete(r.Context(), id)
		if err != nil {
			if err.Error() == "Profile not found" || err.Error() == "sql: no rows in result set" {
				encodeError(w, &models.APIError{
					Status:  http.StatusNotFound,
					Message: "Profile not found",
				})
				return
			}
			encodeError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleSearchProfiles(svc *profile.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		q := query.Get("q")
		if q == "" {
			encodeError(w, &models.APIError{
				Status:  http.StatusBadRequest,
				Message: "Unable to interpret query",
			})
			return
		}

		filter, err := parseSearchQuery(q)
		if err != nil {
			encodeError(w, &models.APIError{
				Status:  http.StatusUnprocessableEntity,
				Message: err.Error(),
			})
			return
		}

		queryOptions, err := parseQueryOptions(query)
		if err != nil {
			encodeError(w, &models.APIError{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid query parameters",
			})
			return
		}

		result, err := svc.List(r.Context(), filter, queryOptions)
		if err != nil {
			encodeError(w, err)
			return
		}

		encode(w, http.StatusOK, models.APIResponse[any]{
			Status: "success",
			Page:   lo.ToPtr(int(queryOptions.Page)),
			Limit:  lo.ToPtr(int(queryOptions.Limit)),
			Total:  lo.ToPtr(int(result.Total)),
			Data:   result.Data,
		})
	}
}
