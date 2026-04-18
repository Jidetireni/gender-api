package postgres

import (
	"strings"

	"github.com/Jidetireni/gender-api/config"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func New(cfg *config.Config) (*PostgresDB, error) {
	db, err := sqlx.Open("postgres", cfg.DBURL)
	if err != nil {
		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	return &PostgresDB{
		DB: db,
	}, nil
}
