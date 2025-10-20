package entity

type Sensor struct {
	ID           string  `json:"Id"`
	Type         string  `json:"type" validate:"oneof=humidity temperature pressure"`
	Alias        string  `json:"alias"`
	Rate         int     `json:"rate"`
	MaxThreshold float32 `json:"maxThreshold"`
	MinThreshold float32 `json:"minThreshold"`
	UpdatedAt    int64   `json:"updatedAt"`
}
