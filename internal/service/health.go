package service

import (
	"context"
	"fmt"

	"github.com/viniosilva/swordhealth-api/internal/repository"
)

//go:generate mockgen -destination=../../mock/health_service_mock.go -package=mock . HealthService
type HealthService interface {
	Health(ctx context.Context) error
}

type healthService struct {
	healthRepository repository.HealthRepository
}

func NewHealthService(healthRepository repository.HealthRepository) HealthService {
	return &healthService{
		healthRepository: healthRepository,
	}
}

func (impl *healthService) Health(ctx context.Context) error {
	err := impl.healthRepository.Health(ctx)
	if err != nil {
		fmt.Println("internal.service.health.health.error: ", err.Error())
	}

	return err
}
