package model

import (
	"testing"
	"time"
)

func TestIngestRequestFromCompactJSON(t *testing.T) {
	body := []byte(`{"id":1,"nom":1,"ena":1,"tIn":"2026-10-22 17:46:13","lmm":99999,"vli":999999.99,"vm3":999999.99,"vto":999999.99,"ind":999999.99,"oud":999999.99,"inm":999999.99,"oum":999999.99,"war":0,"ala":0}`)

	req, err := IngestRequestFromCompactJSON(body)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if req.PlantID != "plant-01" {
		t.Fatalf("plant_id: got %q", req.PlantID)
	}
	if len(req.Readings) != 1 {
		t.Fatalf("readings: got %d", len(req.Readings))
	}
	r := req.Readings[0]
	if r.BoilerNumber != 1 {
		t.Fatalf("boiler_number: got %d", r.BoilerNumber)
	}
	wantTS := time.Date(2026, 10, 22, 17, 46, 13, 0, time.UTC)
	if !r.TS.Equal(wantTS) {
		t.Fatalf("ts: got %v want %v", r.TS, wantTS)
	}
	if r.LevelMM != 99999 || r.VolumeL != 999999.99 || r.VolumeM3 != 999999.99 || r.MassT != 999999.99 {
		t.Fatalf("metrics mismatch: %+v", r)
	}
	if r.IncomeMonth != 999999.99 || r.OutflowMonth != 999999.99 || r.IncomeYear != 999999.99 || r.OutflowYear != 999999.99 {
		t.Fatalf("flow mismatch: %+v", r)
	}
}

func TestLooksLikeCompactPLC(t *testing.T) {
	compact := []byte(`{"id":1,"nom":1,"lmm":3500}`)
	agent := []byte(`{"plant_id":"plant-01","readings":[{"boiler_number":1,"level_mm":1}]}`)
	if !LooksLikeCompactPLC(compact) {
		t.Fatal("expected compact")
	}
	if LooksLikeCompactPLC(agent) {
		t.Fatal("agent envelope must not look compact")
	}
}
