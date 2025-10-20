package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"

	"github.com/AntonioBR9998/go-nats-simulator/gan/config"
	"github.com/AntonioBR9998/go-nats-simulator/utils"
)

type Repository interface {
	MetricRepository
	SensorRepository
}

type repository struct {
	timescaleDbClient *sql.DB
}

func NewRepository(cfg config.Config) Repository {
	timescaleDbClient := NewPostgresClient(cfg.TimescaleDB)

	return &repository{
		timescaleDbClient: timescaleDbClient,
	}
}

func NewPostgresClient(conf utils.PostgreSQLConfig) *sql.DB {
	host := conf.Host
	port := conf.Port
	user := conf.User
	password := conf.Password
	dbName := conf.DBName
	sslMode := conf.SSLMode

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbName, password, sslMode)
	database, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		log.Error("there is an error while connecting to the postgres database ", err)
		panic(err)
	} else {
		log.Info("successfully connected to postgres database!")
		return database
	}

}
