package domain

import (
	"github.com/AntonioBR9998/go-common/validation"
	"github.com/AntonioBR9998/go-nats-simulator/gan/config"
	"github.com/AntonioBR9998/go-nats-simulator/gan/repository"
	"github.com/AntonioBR9998/go-nats-simulator/gan/simulator"
)

type Service interface {
	SensorService
	MetricService
}

type service struct {
	repo      repository.Repository
	conf      config.Config
	validate  *validation.Validator
	simulator *simulator.Manager
}

func NewService(repo repository.Repository, conf config.Config, simulator *simulator.Manager) Service {
	validator, err := validation.NewValidator()
	if err != nil {
		panic(err)
	}

	svc := &service{
		repo:      repo,
		conf:      conf,
		validate:  validator,
		simulator: simulator,
	}

	return svc
}
