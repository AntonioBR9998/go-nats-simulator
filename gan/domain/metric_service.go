package domain

import (
	"context"

	"github.com/AntonioBR9998/go-nats-simulator/gan/domain/entity"
)

type MetricService interface {
	GetMetricsData(ctx context.Context) ([]*entity.Metric, error)
}

func (s *service) GetMetricsData(ctx context.Context) ([]*entity.Metric, error) {
	// Calling repository
	metricsData, err := s.repo.GetMetrics(ctx)
	if err != nil {
		return nil, err
	}

	return metricsData, nil
}
