package repository

const (
	// Sensors
	DEVICE_FIELDS = "id, type, alias, rate, max_threshold, min_threshold, updated_at"

	INSERT_SENSOR = `
		INSERT INTO devices (` + DEVICE_FIELDS + `)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	REPLACE_SENSOR = `
		UPDATE devices
		SET type=$2, alias=$3, rate=$4, max_threshold=$5, min_threshold=$6, updated_at=$7
		WHERE id=$1;`

	DELETE_SENSOR = `
		DELETE FROM devices
		WHERE id=$1;`

	GET_SENSORS = `
        SELECT
			` + DEVICE_FIELDS + `
        FROM devices;`
)
