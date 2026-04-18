package database

import (
	"github.com/Jidetireni/gender-api/config"
	"github.com/Jidetireni/gender-api/internals/pkg/database/postgres"
)

type Database struct {
	PostgresDB *postgres.PostgresDB
}

func New(cfg *config.Config) (*Database, error) {
	pgDB, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Database{
		PostgresDB: pgDB,
	}, nil
}
