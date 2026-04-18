package postgres

import (
	"github.com/Jidetireni/gender-api/config"
	"github.com/jmoiron/sqlx"
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

	return &PostgresDB{
		DB: db,
	}, nil
}
