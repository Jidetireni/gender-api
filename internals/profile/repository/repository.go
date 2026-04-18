package repository

import (
	"context"

	"github.com/Jidetireni/gender-api/internals/profile/handlers/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
	ID        *uuid.UUID
	Gender    *string
	CountryID *string
	AgeGroup  *string
}

func (c *ProfileRepository) buildQuery(filter *ProfileRepositoryFilter, queryType QueryType) (string, []any, error) {
	var builder sq.SelectBuilder
	switch queryType {
	case QueryTypeSelect:
		builder = c.psql.Select("*").From("profiles")

	case QueryTypeCount:
		builder = c.psql.Select("COUNT(*)").From("profiles")
	}

	if filter.ID != nil {
		builder = builder.Where(sq.Eq{"id": filter.ID})
	}
	if filter.Gender != nil {
		builder = builder.Where(sq.ILike{"gender": filter.Gender})
	}
	if filter.CountryID != nil {
		builder = builder.Where(sq.ILike{"country_id": filter.CountryID})
	}
	if filter.AgeGroup != nil {
		builder = builder.Where(sq.ILike{"age_group": filter.AgeGroup})
	}

	return builder.ToSql()
}

func (p *ProfileRepository) Get(ctx context.Context, filter *ProfileRepositoryFilter) (*Profile, error) {
	query, args, err := p.buildQuery(filter, QueryTypeSelect)
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

func (p *ProfileRepository) List(ctx context.Context, filter *ProfileRepositoryFilter) ([]*Profile, error) {
	query, args, err := p.buildQuery(filter, QueryTypeSelect)
	if err != nil {
		return nil, err
	}

	var profiles []*Profile
	err = p.db.SelectContext(ctx, &profiles, query, args...)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (p *ProfileRepository) Upsert(ctx context.Context, profile *Profile) (*Profile, bool, error) {
	builder := p.psql.Insert("profiles").
		Columns("id", "name", "gender", "gender_probability", "sample_size", "age", "age_group", "country_id", "country_probability").
		Values(profile.ID, profile.Name, profile.Gender, profile.GenderProbability, profile.SampleSize, profile.Age, profile.AgeGroup, profile.CountryID, profile.CountryProbability).
		Suffix(`ON CONFLICT (name)
                DO UPDATE SET
                    name = EXCLUDED.name,
                    gender = EXCLUDED.gender,
                    gender_probability = EXCLUDED.gender_probability,
                    sample_size = EXCLUDED.sample_size,
                    age = EXCLUDED.age,
                    age_group = EXCLUDED.age_group,
                    country_id = EXCLUDED.country_id,
                    country_probability = EXCLUDED.country_probability,
                    updated_at = CURRENT_TIMESTAMP
                RETURNING id, name, gender, gender_probability, sample_size, age, age_group, country_id, country_probability, created_at, updated_at, (xmax = 0) AS is_insert`)

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
			&res.SampleSize,
			&res.Age,
			&res.AgeGroup,
			&res.CountryID,
			&res.CountryProbability,
			&res.CreatedAt,
			&res.UpdatedAt,
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
		SampleSize:         int(profile.SampleSize),
		Age:                int(profile.Age),
		AgeGroup:           profile.AgeGroup,
		CountryID:          profile.CountryID,
		CountryProbability: profile.CountryProbability,
		CreatedAt:          profile.CreatedAt,
	}
}
