package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Jidetireni/gender-api/config"
	"github.com/Jidetireni/gender-api/internals/pkg/cache/redis"
	"github.com/Jidetireni/gender-api/internals/profile/handlers/models"
	"github.com/Jidetireni/gender-api/internals/profile/repository"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

var _ ProfileRepository = (*repository.ProfileRepository)(nil)
var _ RedisService = (*redis.Redis)(nil)

type ProfileRepository interface {
	Get(ctx context.Context, filter *repository.ProfileRepositoryFilter) (*repository.Profile, error)
	List(ctx context.Context, filter *repository.ProfileRepositoryFilter) ([]*repository.Profile, error)
	Upsert(ctx context.Context, profile *repository.Profile) (*repository.Profile, bool, error)
	Delete(ctx context.Context, id *uuid.UUID) error
	MapRepositoryToHandlerModel(profile *repository.Profile) *models.Profile
}

type RedisService interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key string) error
}

type Service struct {
	config            config.Config
	httpClient        *http.Client
	profileRepository ProfileRepository
	redisService      RedisService
}

func New(config *config.Config, profileRepository ProfileRepository, redis RedisService) *Service {
	return &Service{
		config: *config,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		profileRepository: profileRepository,
		redisService:      redis,
	}
}

func (s *Service) getGenderResponse(ctx context.Context, name string) (*GenderResponse, error) {
	url := fmt.Sprintf("%s/?name=%s", s.config.GenderizedAPIBaseURL, name)

	genderResponse := &GenderResponse{}
	if err := s.call(ctx, url, genderResponse); err != nil {
		return nil, err
	}

	return genderResponse, nil
}

func (s *Service) getAgeResponse(ctx context.Context, name string) (*AgeResponse, error) {
	url := fmt.Sprintf("%s/?name=%s", s.config.AgifyAPIBaseURL, name)

	ageResponse := &AgeResponse{}
	if err := s.call(ctx, url, ageResponse); err != nil {
		return nil, err
	}

	return ageResponse, nil
}

func (s *Service) getNationalityResponse(ctx context.Context, name string) (*NationalityResponse, error) {
	url := fmt.Sprintf("%s/?name=%s", s.config.NationalizeAPIBaseURL, name)

	nationalityResponse := &NationalityResponse{}
	if err := s.call(ctx, url, nationalityResponse); err != nil {
		return nil, err
	}

	return nationalityResponse, nil
}

func (s *Service) call(ctx context.Context, url string, response any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBody, response); err != nil {
		return err
	}

	return nil
}

func (s *Service) Create(ctx context.Context, name string) (*models.Profile, bool, error) {
	var cachedProfile models.Profile
	err := s.redisService.Get(ctx, RedisProfileNameKey(name), &cachedProfile)
	if err == nil {
		return &cachedProfile, false, nil
	}

	g, gCtx := errgroup.WithContext(ctx)

	var genderResponse *GenderResponse
	var ageResponse *AgeResponse
	var nationalityResponse *NationalityResponse

	g.Go(func() error {
		var err error
		genderResponse, err = s.getGenderResponse(gCtx, name)
		return err
	})

	g.Go(func() error {
		var err error
		ageResponse, err = s.getAgeResponse(gCtx, name)
		return err
	})

	g.Go(func() error {
		var err error
		nationalityResponse, err = s.getNationalityResponse(gCtx, name)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, false, err
	}

	if genderResponse == nil || genderResponse.Count == 0 || genderResponse.Gender == "" {
		return nil, false, &models.APIError{
			Status:  http.StatusBadGateway,
			Message: "Genderize returned an invalid response",
		}
	}

	if ageResponse == nil || ageResponse.Count == 0 || ageResponse.Age == nil {
		return nil, false, &models.APIError{
			Status:  http.StatusBadGateway,
			Message: "Agify returned an invalid response",
		}
	}

	if nationalityResponse == nil || nationalityResponse.Count == 0 || len(nationalityResponse.Country) <= 0 {
		return nil, false, &models.APIError{
			Status:  http.StatusBadGateway,
			Message: "Nationalize returned an invalid response",
		}
	}

	var ageGroup string
	switch {
	case *ageResponse.Age >= 0 && *ageResponse.Age <= 12:
		ageGroup = "child"
	case *ageResponse.Age >= 13 && *ageResponse.Age <= 19:
		ageGroup = "teenager"
	case *ageResponse.Age >= 20 && *ageResponse.Age <= 59:
		ageGroup = "adult"
	case *ageResponse.Age >= 60:
		ageGroup = "senior"
	}

	sampleSize := genderResponse.Count

	id, _ := uuid.NewV7()
	profile, isInsert, err := s.profileRepository.Upsert(ctx, &repository.Profile{
		ID:                 id,
		Name:               name,
		Gender:             genderResponse.Gender,
		GenderProbability:  float64(genderResponse.Probability),
		SampleSize:         int32(sampleSize),
		Age:                int32(*ageResponse.Age),
		AgeGroup:           ageGroup,
		CountryID:          nationalityResponse.Country[0].CountryID,
		CountryProbability: float64(nationalityResponse.Country[0].Probability),
	})
	if err != nil {
		return nil, false, err
	}

	_ = s.redisService.Set(ctx, RedisProfileNameKey(name), profile, RedisProfileExpirationTTL)
	return s.profileRepository.MapRepositoryToHandlerModel(profile), isInsert, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Profile, error) {
	var cachedProfile models.Profile
	err := s.redisService.Get(ctx, RedisProfileIDKey(id.String()), &cachedProfile)
	if err == nil {
		return &cachedProfile, nil
	}

	profile, err := s.profileRepository.Get(ctx, &repository.ProfileRepositoryFilter{
		ID: &id,
	})
	if err != nil {
		return nil, fmt.Errorf("Profile not found")
	}

	_ = s.redisService.Set(ctx, RedisProfileIDKey(id.String()), profile, RedisProfileExpirationTTL)
	return s.profileRepository.MapRepositoryToHandlerModel(profile), nil
}

func (s *Service) List(ctx context.Context, filter *repository.ProfileRepositoryFilter) ([]*models.Profile, error) {
	profiles, err := s.profileRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	mappedProfiles := lo.Map(profiles, func(p *repository.Profile, _ int) *models.Profile {
		return s.profileRepository.MapRepositoryToHandlerModel(p)
	})

	return mappedProfiles, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	profile, err := s.profileRepository.Get(ctx, &repository.ProfileRepositoryFilter{
		ID: &id,
	})
	if err != nil {
		return fmt.Errorf("Profile not found")
	}

	_ = s.redisService.Delete(ctx, RedisProfileIDKey(id.String()))
	_ = s.redisService.Delete(ctx, RedisProfileNameKey(profile.Name))

	return s.profileRepository.Delete(ctx, &id)
}
