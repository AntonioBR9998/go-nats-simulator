package entity

type Metric struct {
	SensorID  string  `json:"sensorId"`
	Value     float32 `json:"value"`
	Unit      string  `json:"unit"`
	Timestamp int64   `json:"timestamp"`
}
