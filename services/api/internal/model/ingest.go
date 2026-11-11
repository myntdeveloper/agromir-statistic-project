package model

import "time"

type IngestRequest struct {
	PlantID  string          `json:"plant_id" binding:"required"`
	SentAt   time.Time       `json:"sent_at"`
	Readings []IngestReading `json:"readings" binding:"required,min=1,dive"`
}

type IngestReading struct {
	BoilerNumber int       `json:"boiler_number" binding:"required,min=1"`
	TS           time.Time `json:"ts"`
	LevelMM      float64   `json:"level_mm"`
	Density      float64   `json:"density"`
	VolumeM3     float64   `json:"volume_m3"`
	VolumeL      float64   `json:"volume_l"`
	MassT        float64   `json:"mass_t"`
	IncomeMonth  float64   `json:"income_month"`
	OutflowMonth float64   `json:"outflow_month"`
	IncomeYear   float64   `json:"income_year"`
	OutflowYear  float64   `json:"outflow_year"`
}

type IngestResponse struct {
	Accepted int `json:"accepted"`
}

type PlantStatus struct {
	PlantID        string    `json:"plant_id"`
	LastAcceptedAt time.Time `json:"last_accepted_at"`
	LastCount      int       `json:"last_count"`
}
