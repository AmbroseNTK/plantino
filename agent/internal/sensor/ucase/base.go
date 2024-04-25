package ucase

import (
	"context"

	"github.com/ambrosentk/plantino/agent/domain"
)

type SensorDataUseCase struct {
	repo domain.SensorDataRepository
}

func (s SensorDataUseCase) Create(ctx context.Context, data *domain.SensorData) error {
	return s.repo.Create(ctx, data)
}

func NewSensorDataUseCase(repo domain.SensorDataRepository) *SensorDataUseCase {
	return &SensorDataUseCase{repo: repo}
}
