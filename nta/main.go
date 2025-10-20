// nta could be a complete service with hexagonal structure instead of a simple script
// for scalability. Simple script has been chosen because of quick development

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

const (
	NATS_URL = "nats://nats:4222"
	HOST     = "timescale-db"
	PORT     = 5432
	USER     = "admin"
	PASS     = "admin"
	DBNAME   = "sensors"
	SSLMODE  = "disable"
)

type Sample struct {
	SensorID  string
	Value     float32
	Unit      string
	Timestamp int64
}

func main() {
	log.Print("starting NTA (NATS TimescaleDB Adapter)")

	// Connecting TimescaleDB
	log.Print("connecting to timescaleDB")

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", HOST, PORT, USER, DBNAME, PASS, SSLMODE)
	db, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		log.Printf("there is an error while connecting to the database: %v", err)
	}

	defer db.Close()

	// Connecting NATS
	log.Print("connecting to NATS")

	natsClient, err := nats.Connect(NATS_URL)
	if err != nil {
		log.Printf("error connecting to NATS: %v", err)
	}
	defer natsClient.Close()

	// Subscribing to sensors topic
	log.Print("subscribing to sensors topic")
	sub, err := natsClient.Subscribe("sensors", func(msg *nats.Msg) {
		var event Sample
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("event is not processable: %v", err)
			return
		}

		// Inserting event in TimescaleDB
		query := `INSERT INTO metrics (sensor_id, value, unit, timestamp) VALUES ($1, $2, $3, $4)`

		if _, err := db.Exec(
			query,
			event.SensorID,
			event.Value,
			event.Unit,
			event.Timestamp); err != nil {
			log.Printf("error writing in database: %v", err)
			return
		}
	})
	if err != nil {
		log.Printf("error subscribing sensors topic: %v", err)
	}

	// Waiting for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Consumidor detenido")

	defer sub.Unsubscribe()
}
