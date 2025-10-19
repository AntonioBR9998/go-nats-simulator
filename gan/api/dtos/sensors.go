package dtos

import (
	"github.com/AntonioBR9998/go-nats-simulator/gan/domain/entity"
)

type SensorBaseRequest struct {
	Body SensorBody `contentType:"application/json"`
}

type SensorBody struct {
	ID           string  `json:"Id"`
	Type         string  `json:"type"`
	Alias        string  `json:"alias"`
	Rate         int     `json:"rate"`
	MaxThreshold float32 `json:"maxThreshold"`
	MinThreshold float32 `json:"minThreshold"`
}

type SensorResponse struct {
	ID           string  `json:"Id"`
	Type         string  `json:"type"`
	Alias        string  `json:"alias"`
	Rate         int     `json:"rate"`
	MaxThreshold float32 `json:"maxThreshold"`
	MinThreshold float32 `json:"minThreshold"`
	UpdatedAt    int64   `json:"updatedAt"`
}

func ToSensorResponseDto(res *entity.Sensor) SensorResponse {
	return SensorResponse{
		ID:           res.ID,
		Type:         res.Type,
		Alias:        res.Alias,
		Rate:         res.Rate,
		MaxThreshold: res.MaxThreshold,
		MinThreshold: res.MinThreshold,
		UpdatedAt:    res.UpdatedAt,
	}
}

func ValidateSensorType(typ string) bool {
	switch typ {
	case "humidity", "temperature", "pressure":
		return true
	default:
		return false
	}
}
