package domain

import (
	"context"
	"time"

	"github.com/AntonioBR9998/go-nats-simulator/errors"
	"github.com/AntonioBR9998/go-nats-simulator/gan/domain/entity"
	log "github.com/sirupsen/logrus"
)

type SensorService interface {
	CreateSensor(ctx context.Context, id string, typ string, alias string, rate int,
		maxTh float32, minTh float32) (*entity.Sensor, error)
	ModifySensor(ctx context.Context, id string, typ string, alias string, rate int,
		maxTh float32, minTh float32) (*entity.Sensor, error)
	GetSensors(ctx context.Context) ([]*entity.Sensor, error)
	DeleteSensor(ctx context.Context, id string) error
}

func (s *service) CreateSensor(ctx context.Context,
	id string,
	typ string,
	alias string,
	rate int,
	maxTh float32,
	minTh float32,
) (*entity.Sensor, error) {
	errVars := map[string]any{"id": id, "alias": alias}

	// Validating param id
	err := s.validate.Var(id, "uuid_rfc4122")
	if err != nil {
		log.Errorln("validation error: ", err)
		return nil, errors.TrackErrorVar(err, errVars)
	}

	// Validating param alias
	err = s.validate.Var(alias, "128_character_name")
	if err != nil {
		log.Errorln("validation error: ", err)
		return nil, errors.TrackErrorVar(err, errVars)
	}

	// Adding sensor in database
	updatedAt := time.Now().Unix()

	sensor := &entity.Sensor{
		ID:           id,
		Type:         typ,
		Alias:        alias,
		Rate:         rate,
		MaxThreshold: maxTh,
		MinThreshold: minTh,
		UpdatedAt:    updatedAt,
	}

	err = s.repo.CreateSensor(ctx, sensor)
	if err != nil {
		return nil, err
	}

	// Adding sensor to simulator
	go s.simulator.Start(sensor.ID, sensor.Type, sensor.Rate)

	return sensor, nil
}

func (s *service) ModifySensor(ctx context.Context,
	id string,
	typ string,
	alias string,
	rate int,
	maxTh float32,
	minTh float32,
) (*entity.Sensor, error) {
	errVars := map[string]any{"id": id, "alias": alias}

	// Validating param id
	err := s.validate.Var(id, "uuid_rfc4122")
	if err != nil {
		log.Errorln("validation error: ", err)
		return nil, errors.TrackErrorVar(err, errVars)
	}

	// Validating param alias
	err = s.validate.Var(alias, "128_character_name")
	if err != nil {
		log.Errorln("validation error: ", err)
		return nil, errors.TrackErrorVar(err, errVars)
	}

	// Updating sensor in database
	updatedAt := time.Now().Unix()

	sensor := &entity.Sensor{
		ID:           id,
		Type:         typ,
		Alias:        alias,
		Rate:         rate,
		MaxThreshold: maxTh,
		MinThreshold: minTh,
		UpdatedAt:    updatedAt,
	}

	err = s.repo.ModifySensor(ctx, sensor)
	if err != nil {
		return nil, err
	}

	// Replacing sensor in simulator
	go s.simulator.Start(sensor.ID, sensor.Type, sensor.Rate)

	return sensor, nil
}

func (s *service) GetSensors(ctx context.Context) ([]*entity.Sensor, error) {
	// Calling repository
	sensorList, err := s.repo.GetSensors(ctx)
	if err != nil {
		return nil, err
	}

	return sensorList, nil
}

func (s *service) DeleteSensor(ctx context.Context, id string) error {
	errVars := map[string]any{"id": id}

	// Validating param id
	err := s.validate.Var(id, "uuid_rfc4122")
	if err != nil {
		log.Errorln("validation error: ", err)
		return errors.TrackErrorVar(err, errVars)
	}

	// Deleting sensor in database
	err = s.repo.DeleteSensor(ctx, id)
	if err != nil {
		return err
	}

	// Deleting sensor in simulator
	s.simulator.Stop(id)

	return nil
}
