package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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

// plcCompact is the short single-tank JSON from S7/WinCC, e.g.
// {"id":1,"nom":1,"ena":1,"tIn":"2026-10-22 17:46:13","lmm":99999,"vli":...,"vm3":...,"vto":...,"ind":...,"oud":...,"inm":...,"oum":...,"war":0,"ala":0}
type plcCompact struct {
	ID  int     `json:"id"`
	Nom int     `json:"nom"`
	Ena int     `json:"ena"`
	TIn string  `json:"tIn"`
	Lmm float64 `json:"lmm"`
	Vli float64 `json:"vli"`
	Vm3 float64 `json:"vm3"`
	Vto float64 `json:"vto"`
	Ind float64 `json:"ind"`
	Oud float64 `json:"oud"`
	Inm float64 `json:"inm"`
	Oum float64 `json:"oum"`
	War int     `json:"war"`
	Ala int     `json:"ala"`
}

// LooksLikeCompactPLC reports the short PLC body (no plant_id/readings wrapper).
func LooksLikeCompactPLC(body []byte) bool {
	s := string(body)
	if strings.Contains(s, `"plant_id"`) || strings.Contains(s, `"readings"`) {
		return false
	}
	return strings.Contains(s, `"lmm"`) || strings.Contains(s, `"nom"`) || strings.Contains(s, `"tIn"`)
}

// IngestRequestFromCompactJSON maps one compact PLC object into IngestRequest.
// id → plant-%02d, nom → boiler_number; lmm/vli/vm3/vto/inm/oum/ind/oud → reading fields.
func IngestRequestFromCompactJSON(body []byte) (IngestRequest, error) {
	var c plcCompact
	if err := json.Unmarshal(body, &c); err != nil {
		return IngestRequest{}, fmt.Errorf("compact plc json: %w", err)
	}

	boiler := c.Nom
	if boiler < 1 {
		boiler = c.ID
	}
	if boiler < 1 {
		return IngestRequest{}, fmt.Errorf("compact plc: nom (boiler number) is required")
	}

	plantID := "plant-01"
	if c.ID > 0 {
		plantID = fmt.Sprintf("plant-%02d", c.ID)
	}

	ts, err := parsePLCTime(c.TIn)
	if err != nil {
		return IngestRequest{}, err
	}

	reading := IngestReading{
		BoilerNumber: boiler,
		TS:           ts,
		LevelMM:      c.Lmm,
		VolumeL:      c.Vli,
		VolumeM3:     c.Vm3,
		MassT:        c.Vto,
		IncomeMonth:  c.Inm,
		OutflowMonth: c.Oum,
		IncomeYear:   c.Ind,
		OutflowYear:  c.Oud,
	}

	return IngestRequest{
		PlantID:  plantID,
		SentAt:   ts,
		Readings: []IngestReading{reading},
	}, nil
}

func parsePLCTime(raw string) (time.Time, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return time.Time{}, nil
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("compact plc: invalid tIn %q", raw)
}

type IngestResponse struct {
	Accepted int `json:"accepted"`
}

type PlantStatus struct {
	PlantID        string    `json:"plant_id"`
	LastAcceptedAt time.Time `json:"last_accepted_at"`
	LastCount      int       `json:"last_count"`
}
