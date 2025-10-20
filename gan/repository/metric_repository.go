package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AntonioBR9998/go-nats-simulator/errors"
	"github.com/AntonioBR9998/go-nats-simulator/gan/domain/entity"
	"github.com/AntonioBR9998/go-nats-simulator/humamw"
	"github.com/AntonioBR9998/go-nats-simulator/humamw/sql"
	log "github.com/sirupsen/logrus"
)

type MetricRepository interface {
	GetMetrics(ctx context.Context) ([]*entity.Metric, error)
}

// Allowed fields to filter by in /GET metrics
var getMetricsWhereDef = map[string]string{
	"sensorId":  "sensor_id",
	"timestamp": "timestamp",
}

// Allowed fields to order by in /GET metrics
var getMetricsOrderDef = map[string]string{
	"value":     "value",
	"timestamp": "timestamp",
}

func (r *repository) GetMetrics(ctx context.Context) ([]*entity.Metric, error) {
	log.Debug("getting metrics in repository")

	queryTemplate := GET_METRICS
	args := []any{}

	// Getting filters
	filter, hasFilter := humamw.GetFilter(ctx)
	if hasFilter {
		var err error
		queryTemplate, args, err = sql.AddFilterToQuery(filter, getMetricsWhereDef, getMetricsOrderDef, queryTemplate, args)
		if err != nil {
			return nil, errors.TrackError(err)
		}
	}

	// Getting total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS subquery", queryTemplate)
	var total int
	prev := time.Now()
	err := r.timescaleDbClient.QueryRow(countQuery, args...).Scan(&total)
	log.Infof("Time counting metrics: %v", time.Since(prev))

	if err != nil {
		log.Errorf("Error while counting metrics: %v \n query: %s", err, countQuery)
		return nil, errors.TrackError(err)
	}

	// Pagination should has default values (see middleware to setup)
	pagination, hasPagination := humamw.GetPagination(ctx)
	// If has pagination
	if hasPagination {
		// Update template to use pagination
		argsLen := len(args)
		paginatedStr := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argsLen+1, argsLen+2)
		queryTemplate = strings.Join([]string{queryTemplate, paginatedStr, ";"}, "")
		// Add pagination limit and offset
		args = append(args, pagination.Limit, pagination.Offset)
	}

	rows, err := r.timescaleDbClient.Query(queryTemplate, args...)
	if err != nil {
		log.Errorf("Error executing query: %s \n error: %v", queryTemplate, err)
		return nil, errors.TrackError(err)
	}
	defer rows.Close()

	// Reading rows
	var metricsData = []*entity.Metric{}
	for rows.Next() {
		var metric entity.Metric

		if err := rows.Scan(&metric.SensorID, &metric.Value, &metric.Unit, &metric.Timestamp); err != nil {
			log.Errorln("Error scanning metrics table rows:", err)
			return nil, errors.TrackError(err)
		}

		metricsData = append(metricsData, &metric)
	}

	if cb, ok := humamw.GetSetHeaderCallback(ctx, "Total"); ok {
		cb("Total", strconv.Itoa(total))
	}

	return metricsData, nil
}
