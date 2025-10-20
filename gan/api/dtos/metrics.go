package dtos

import (
	"github.com/AntonioBR9998/go-nats-simulator/gan/domain/entity"
)

type MetricResponse struct {
	SensorID  string  `json:"sensorId"`
	Value     float32 `json:"value"`
	Unit      string  `json:"unit"`
	Timestamp int64   `json:"timestamp"`
}

func ToMetricResponseDto(res *entity.Metric) *MetricResponse {
	return &MetricResponse{
		SensorID:  res.SensorID,
		Value:     res.Value,
		Unit:      res.Unit,
		Timestamp: res.Timestamp,
	}
}
