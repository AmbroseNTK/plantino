package repo

import (
	"context"

	"github.com/ambrosentk/plantino/agent/domain"
	"gorm.io/gorm"
)

type SensorDataRepository struct {
	db *gorm.DB
}

func (s SensorDataRepository) Create(ctx context.Context, data *domain.SensorData) error {
	tx := s.db.Create(data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func NewSensorDataRepository(db *gorm.DB) *SensorDataRepository {
	err := db.AutoMigrate(&domain.SensorData{})
	if err != nil {
		panic(err)
	}
	return &SensorDataRepository{db: db}
}
