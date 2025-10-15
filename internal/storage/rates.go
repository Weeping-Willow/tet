package storage

import "github.com/jmoiron/sqlx"

type RateRepository struct {
	db *sqlx.DB
}

func NewRateRepository(db *sqlx.DB) *RateRepository {
	return &RateRepository{db: db}
}
