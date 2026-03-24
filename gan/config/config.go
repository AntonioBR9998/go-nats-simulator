package config

import (
	"github.com/AntonioBR9998/go-common/config"
)

// Config represents the service configuration
type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	Nats        NatsConfig              `json:"nats"`
	TimescaleDB config.PostgreSQLConfig `json:"timescaleDB"`
	ServerName  string                  `json:"serverName"`
}

type NatsConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}
