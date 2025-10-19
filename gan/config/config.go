package config

import (
	"github.com/AntonioBR9998/go-nats-simulator/utils"
)

// Config represents the service configuration
type Config struct {
	utils.BaseConfig `mapstructure:",squash"`

	Nats        NatsConfig             `json:"nats"`
	TimescaleDB utils.PostgreSQLConfig `json:"timescaleDB"`
	ServerName  string                 `json:"serverName"`
}

type NatsConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
