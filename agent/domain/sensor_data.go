package domain

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type SensorData struct {
	gorm.Model
	Temperature     float64 `json:"temperature"`
	Humidity        float64 `json:"humidity"`
}

type SensorDataRepository interface {
	Create(ctx context.Context, data *SensorData) error
}

type SensorDataUsecase interface {
	Create(ctx context.Context, data *SensorData) error
}

var (
	ErrInvalidSensorData = errors.New("invalid sensor data")
)

func NewSensorDataFromRawData(rawData []byte) (*SensorData, error) {
	strData := string(rawData)
	
	jsonData := &SensorData{}

	_ = json.Unmarshal([]byte(strData), jsonData)

	return jsonData, nil
}
