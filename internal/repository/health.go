package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=../../mock/health_repository_mock.go -package=mock . HealthRepository
type HealthRepository interface {
	Health(ctx context.Context) error
}

type healthRepository struct {
	db *sqlx.DB
}

func NewHealthRepository(db *sqlx.DB) HealthRepository {
	return &healthRepository{
		db: db,
	}
}

func (impl *healthRepository) Health(ctx context.Context) error {
	return impl.db.PingContext(ctx)
}
