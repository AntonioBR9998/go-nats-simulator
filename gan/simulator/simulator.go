package simulator

import (
	"encoding/json"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/AntonioBR9998/go-nats-simulator/gan/domain/entity"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

// It manages actives sensors
type Manager struct {
	natsClient *nats.Conn
	simulators map[string]chan struct{}
	mu         sync.Mutex
}

func NewManager(natsClient *nats.Conn) *Manager {
	return &Manager{
		natsClient: natsClient,
		simulators: make(map[string]chan struct{}),
	}
}

// This function initializes a sensor simulator
func (m *Manager) Start(id string, typ string, rate int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Checking if a sensor with this ID exists and deleting it
	if _, exists := m.simulators[id]; exists {
		log.Warnf("replacing sensor with ID: %s", id)
		m.Stop(id)
	}

	stopCh := make(chan struct{})
	m.simulators[id] = stopCh

	go m.run(id, typ, rate, stopCh)
	log.Infof("new sensor running with ID: %s", id)
}

// This function deletes a sensor
func (m *Manager) Stop(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stopCh, exists := m.simulators[id]
	if !exists {
		return
	}
	close(stopCh)
	delete(m.simulators, id)
	log.Infof("sensor with ID %s has been deleted", id)
}

// Sensor go rutine
func (m *Manager) run(id, typ string, rate int, stopCh <-chan struct{}) {
	subject := "sensors"

	for {
		select {
		case <-stopCh:
			log.Infof("sensor %s stopped", id)
			return
		default:
			value, unit := m.generateValue(typ)

			event := entity.Metric{
				SensorID:  id,
				Value:     value,
				Unit:      unit,
				Timestamp: time.Now().Unix(),
			}

			data, _ := json.Marshal(event)
			if err := m.natsClient.Publish(subject, data); err != nil {
				log.Errorf("error sending data to NATS: %v", err)
			}

			time.Sleep(time.Duration(rate))
		}
	}
}

// This function generates a random value for every type of sensor
func (m *Manager) generateValue(typ string) (float32, string) {
	switch typ {
	case "temperature":
		// Range of temperature: [-30, 60)
		value := rand.Float32()*(float32(60)-float32(-30)) + float32(-30)
		return value, "celsius"
	case "pressure":
		// Range of pressure: [700, 820)
		value := rand.Float32()*(float32(820)-float32(700)) + float32(700)
		return value, "mmHg"
	case "humidity":
		// Range of humidity: [-10, 120)
		value := rand.Float32()*(float32(120)-float32(-10)) + float32(-10)
		return value, "percentage"
	default:
		return float32(0), ""
	}
}
