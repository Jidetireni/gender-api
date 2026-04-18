package handlers

import (
	"net/http"

	"github.com/Jidetireni/gender-api/internals/profile"
	"github.com/Jidetireni/gender-api/internals/profile/handlers/models"
	"github.com/Jidetireni/gender-api/internals/profile/repository"
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
		filter := &repository.ProfileRepositoryFilter{}

		if gender := r.URL.Query().Get("gender"); gender != "" {
			filter.Gender = &gender
		}
		if countryID := r.URL.Query().Get("country_id"); countryID != "" {
			filter.CountryID = &countryID
		}
		if ageGroup := r.URL.Query().Get("age_group"); ageGroup != "" {
			filter.AgeGroup = &ageGroup
		}

		profiles, err := svc.List(r.Context(), filter)
		if err != nil {
			encodeError(w, err)
			return
		}

		shortProfiles := lo.Map(profiles, func(p *models.Profile, _ int) *models.ProfileShort {
			return &models.ProfileShort{
				ID:        p.ID,
				Name:      p.Name,
				Gender:    p.Gender,
				Age:       p.Age,
				AgeGroup:  p.AgeGroup,
				CountryID: p.CountryID,
			}
		})

		encode(w, http.StatusOK, models.APIResponse[[]*models.ProfileShort]{
			Status: "success",
			Count:  lo.ToPtr(len(shortProfiles)),
			Data:   shortProfiles,
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
