#!/bin/bash
set -e

echo "setting up architecture"
docker compose up -d

echo "waiting for timescaleDB be ready"
until docker exec -it timescale-db pg_isready -U admin > /dev/null 2>&1; do
  sleep 3
done

echo "timescaleDB is ready. Creating timestaleDB extension"
docker exec -i timescale-db psql -U admin -d sensors -c "CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;"

echo "creating devices table"
docker exec -i timescale-db psql -U admin -d sensors -c "
CREATE TYPE device_type AS ENUM ('humidity', 'temperature', 'pressure');"
docker exec -i timescale-db psql -U admin -d sensors -c "
CREATE TABLE IF NOT EXISTS devices (
    id VARCHAR PRIMARY KEY,
    type device_type NOT NULL,
    alias INTEGER NOT NULL,
    rate INTEGER NOT NULL,
    max_threshold INTEGER NOT NULL,
    min_threshold INTEGER NOT NULL,
    updated_at BIGINT NOT NULL
);"

echo "creating metrics table"
docker exec -i timescale-db psql -U admin -d sensors -c "
CREATE TABLE IF NOT EXISTS metrics (
    sensor_id VARCHAR NOT NULL,
    value REAL NOT NULL,
    unit VARCHAR NOT NULL,
    timestamp BIGINT NOT NULL
);"
docker exec -i timescale-db psql -U admin -d sensors -c "
SELECT create_hypertable('metrics', 'timestamp', if_not_exists => TRUE);"

echo "the architecture is ready"
