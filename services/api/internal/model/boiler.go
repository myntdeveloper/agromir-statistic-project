package model

import "time"

type WarehouseSummary struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	TotalVolumeM3 float64 `json:"total_volume_m3"`
}

type BoilerListItem struct {
	ID       int     `json:"id"`
	Number   int     `json:"number"`
	Position int     `json:"position"`
	LevelMM  float64 `json:"level_mm"`
	VolumeM3 float64 `json:"volume_m3"`
	IsEmpty  bool    `json:"is_empty"`
}

type BoilerDetail struct {
	ID            int        `json:"id"`
	Number        int        `json:"number"`
	Position      int        `json:"position"`
	WarehouseID   int        `json:"warehouse_id"`
	WarehouseName string     `json:"warehouse_name"`
	DiameterMM    int        `json:"diameter_mm"`
	HeightMaxMM   int        `json:"height_max_mm"`
	UpdatedAt     *time.Time `json:"updated_at"`
	LevelMM       float64    `json:"level_mm"`
	Density       float64    `json:"density"`
	VolumeM3      float64    `json:"volume_m3"`
	VolumeL       float64    `json:"volume_l"`
	MassT         float64    `json:"mass_t"`
	IncomeMonth   float64    `json:"income_month"`
	OutflowMonth  float64    `json:"outflow_month"`
	IncomeYear    float64    `json:"income_year"`
	OutflowYear   float64    `json:"outflow_year"`
	IsEmpty       bool       `json:"is_empty"`
}

type TrendPoint struct {
	TS       time.Time `json:"ts"`
	LevelMM  float64   `json:"level_mm"`
	VolumeM3 float64   `json:"volume_m3"`
	MassT    float64   `json:"mass_t"`
}

type BoilerTrend struct {
	BoilerID int          `json:"boiler_id"`
	Points   []TrendPoint `json:"points"`
}
