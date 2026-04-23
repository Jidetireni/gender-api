package repository

import (
	"context"

	"github.com/Jidetireni/gender-api/internals/profile/handlers/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type ProfileRepository struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type ProfileRepositoryFilter struct {
	ID                    *uuid.UUID
	Gender                *string
	CountryID             *string
	AgeGroup              *string
	MinAge                *int
	MaxAge                *int
	MinGenderProbability  *float64
	MinCountryProbability *float64
}

func (c *ProfileRepository) buildQuery(filter *ProfileRepositoryFilter, queryOptions QueryOptions) (string, []any, error) {
	queryType := lo.FromPtrOr(queryOptions.Type, QueryTypeSelect)

	var builder sq.SelectBuilder
	switch queryType {
	case QueryTypeSelect:
		builder = c.psql.Select("*").From("profiles")
	case QueryTypeCount:
		builder = c.psql.Select("COUNT(*)").From("profiles")
	}

	// Exact filters
	if filter.ID != nil {
		builder = builder.Where(sq.Eq{"id": filter.ID})
	}
	if filter.Gender != nil {
		builder = builder.Where(sq.Eq{"gender": filter.Gender})
	}
	if filter.CountryID != nil {
		builder = builder.Where(sq.Eq{"country_id": filter.CountryID})
	}
	if filter.AgeGroup != nil {
		builder = builder.Where(sq.Eq{"age_group": filter.AgeGroup})
	}

	if filter.MinAge != nil {
		builder = builder.Where(sq.GtOrEq{"age": filter.MinAge})
	}
	if filter.MaxAge != nil {
		builder = builder.Where(sq.LtOrEq{"age": filter.MaxAge})
	}
	if filter.MinGenderProbability != nil {
		builder = builder.Where(sq.GtOrEq{"gender_probability": filter.MinGenderProbability})
	}
	if filter.MinCountryProbability != nil {
		builder = builder.Where(sq.GtOrEq{"country_probability": filter.MinCountryProbability})
	}

	if queryType == QueryTypeSelect {
		var err error
		builder, err = ApplyPagination(builder, queryOptions)
		if err != nil {
			return "", nil, err
		}
	}

	return builder.ToSql()
}

func (p *ProfileRepository) Get(ctx context.Context, filter *ProfileRepositoryFilter) (*Profile, error) {
	query, args, err := p.buildQuery(filter, QueryOptions{Type: lo.ToPtr(QueryTypeSelect)})
	if err != nil {
		return nil, err
	}

	var profile Profile
	err = p.db.GetContext(ctx, &profile, query, args...)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (p *ProfileRepository) Count(ctx context.Context, filter *ProfileRepositoryFilter) (int64, error) {
	query, args, err := p.buildQuery(filter, QueryOptions{Type: lo.ToPtr(QueryTypeCount)})
	if err != nil {
		return 0, err
	}

	var total int64
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func (p *ProfileRepository) List(ctx context.Context, filter *ProfileRepositoryFilter, queryOptions QueryOptions) (*ListResult[Profile], error) {
	query, args, err := p.buildQuery(filter, queryOptions)
	if err != nil {
		return nil, err
	}

	var profiles []*Profile
	err = p.db.SelectContext(ctx, &profiles, query, args...)
	if err != nil {
		return nil, err
	}

	return &ListResult[Profile]{
		Items: profiles,
		Page:  queryOptions.Page,
		Limit: queryOptions.Limit,
	}, nil
}

func (p *ProfileRepository) Upsert(ctx context.Context, profile *Profile) (*Profile, bool, error) {
	builder := p.psql.Insert("profiles").
		Columns(
			"id", "name", "gender", "gender_probability",
			"age", "age_group", "country_id", "country_name", "country_probability",
		).
		Values(
			profile.ID, profile.Name, profile.Gender, profile.GenderProbability,
			profile.Age, profile.AgeGroup, profile.CountryID, profile.CountryName, profile.CountryProbability,
		).
		Suffix(`ON CONFLICT (name) DO NOTHING
			RETURNING id, name, gender, gender_probability,
				age, age_group, country_id, country_name, country_probability,
				created_at, (xmax = 0) AS is_insert`)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, false, err
	}

	var isInsert bool
	var res Profile
	err = p.db.QueryRowContext(ctx, query, args...).
		Scan(
			&res.ID,
			&res.Name,
			&res.Gender,
			&res.GenderProbability,
			&res.Age,
			&res.AgeGroup,
			&res.CountryID,
			&res.CountryName,
			&res.CountryProbability,
			&res.CreatedAt,
			&isInsert,
		)
	if err != nil {
		return nil, false, err
	}

	return &res, isInsert, nil
}

func (p *ProfileRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	query, args, err := p.psql.Delete("profiles").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *ProfileRepository) MapRepositoryToHandlerModel(profile *Profile) *models.Profile {
	return &models.Profile{
		ID:                 profile.ID,
		Name:               profile.Name,
		Gender:             profile.Gender,
		GenderProbability:  profile.GenderProbability,
		Age:                int(profile.Age),
		AgeGroup:           profile.AgeGroup,
		CountryID:          profile.CountryID,
		CountryName:        profile.CountryName,
		CountryProbability: profile.CountryProbability,
		CreatedAt:          profile.CreatedAt,
	}
}
