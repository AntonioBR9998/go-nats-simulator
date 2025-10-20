// Sensors CRUD in TimescaleDB

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

const SENSOR_RESOURCE_TYPE = "sensor"

type SensorRepository interface {
	CreateSensor(ctx context.Context, sensor *entity.Sensor) error
	ModifySensor(ctx context.Context, sensor *entity.Sensor) error
	GetSensors(ctx context.Context) ([]*entity.Sensor, error)
	DeleteSensor(ctx context.Context, id string) error
}

func (r *repository) CreateSensor(ctx context.Context, sensor *entity.Sensor) error {
	log.Debugf("writing in repository devices table a new sensor with ID: %s", sensor.ID)

	errVars := map[string]any{"id": sensor.ID, "alias": sensor.Alias}

	// Writing in TimescaleDB a new sensor
	_, err := r.timescaleDbClient.Exec(
		INSERT_SENSOR,
		sensor.ID,
		sensor.Type,
		sensor.Alias,
		sensor.Rate,
		sensor.MaxThreshold,
		sensor.MinThreshold,
		sensor.UpdatedAt,
	)

	if err != nil {
		err := errors.WrapPostgresErrorCode(err, SENSOR_RESOURCE_TYPE, sensor.ID)
		return errors.TrackErrorVar(err, errVars)
	}

	return nil
}

func (r *repository) ModifySensor(ctx context.Context, sensor *entity.Sensor) error {
	log.Debugf("updating in repository devices table the sensor with ID: %s", sensor.ID)

	errVars := map[string]any{"id": sensor.ID, "alias": sensor.Alias}

	// Updating in TimescaleDB
	_, err := r.timescaleDbClient.Exec(
		REPLACE_SENSOR,
		sensor.ID,
		sensor.Type,
		sensor.Alias,
		sensor.Rate,
		sensor.MaxThreshold,
		sensor.MinThreshold,
		sensor.UpdatedAt,
	)

	if err != nil {
		err := errors.WrapPostgresErrorCode(err, SENSOR_RESOURCE_TYPE, sensor.ID)
		return errors.TrackErrorVar(err, errVars)
	}

	return nil
}

// Allowed fields to filter by in /GET sensors
var getSensorsWhereDef = map[string]string{
	"id":    "id",
	"type":  "type",
	"alias": "alias",
}

// Allowed fields to order by in /GET sensors
var getSensorsOrderDef = map[string]string{
	"id":        "id",
	"type":      "type",
	"alias":     "alias",
	"updatedAt": "updated_at",
}

func (r *repository) GetSensors(ctx context.Context) ([]*entity.Sensor, error) {
	log.Debug("getting sensors in repository")

	queryTemplate := GET_SENSORS
	args := []any{}

	// Getting filters
	filter, hasFilter := humamw.GetFilter(ctx)
	if hasFilter {
		var err error
		queryTemplate, args, err = sql.AddFilterToQuery(filter, getSensorsWhereDef, getSensorsOrderDef, queryTemplate, args)
		if err != nil {
			return nil, errors.TrackError(err)
		}
	}

	// Getting total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS subquery", queryTemplate)
	var total int
	prev := time.Now()
	err := r.timescaleDbClient.QueryRow(countQuery, args...).Scan(&total)
	log.Infof("Time counting sensors: %v", time.Since(prev))

	if err != nil {
		log.Errorf("Error while counting sensors: %v \n query: %s", err, countQuery)
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
	var sensorsData = []*entity.Sensor{}
	for rows.Next() {
		var sensor entity.Sensor

		if err := rows.Scan(&sensor.ID, &sensor.Type, &sensor.Alias, &sensor.Rate,
			&sensor.MaxThreshold, &sensor.MinThreshold, &sensor.UpdatedAt); err != nil {
			log.Errorln("Error scanning devices table rows:", err)
			return nil, errors.TrackError(err)
		}

		sensorsData = append(sensorsData, &sensor)
	}

	if cb, ok := humamw.GetSetHeaderCallback(ctx, "Total"); ok {
		cb("Total", strconv.Itoa(total))
	}

	return sensorsData, nil

}

func (r *repository) DeleteSensor(ctx context.Context, id string) error {
	log.Debugf("deleting in repository devices table the sensor with ID: %s", id)

	errVars := map[string]any{"id": id}

	// Deleting in TimescaleDB
	_, err := r.timescaleDbClient.Exec(
		DELETE_SENSOR,
		id,
	)

	if err != nil {
		err := errors.WrapPostgresErrorCode(err, SENSOR_RESOURCE_TYPE, id)
		return errors.TrackErrorVar(err, errVars)
	}

	return nil
}
