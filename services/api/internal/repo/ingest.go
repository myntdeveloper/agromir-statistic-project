package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/myntdeveloper/agromir_boiler_server/internal/model"
)

type IngestRepo struct {
	db *sql.DB
}

func NewIngestRepo(db *sql.DB) *IngestRepo {
	return &IngestRepo{db: db}
}

func (r *IngestRepo) MapBoilerNumbers(ctx context.Context, numbers []int) (map[int]int, error) {
	out := make(map[int]int, len(numbers))
	if len(numbers) == 0 {
		return out, nil
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT number, id
		FROM boilers
		WHERE number = ANY($1)
	`, pq.Array(numbers))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var number, id int
		if err := rows.Scan(&number, &id); err != nil {
			return nil, err
		}
		out[number] = id
	}
	return out, rows.Err()
}

func (r *IngestRepo) SavePacket(ctx context.Context, req model.IngestRequest, boilerIDs map[int]int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	snapStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO boiler_readings (
			boiler_id, ts, level_mm, density, volume_m3, volume_l, mass_t,
			income_month, outflow_month, income_year, outflow_year
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (boiler_id) DO UPDATE SET
			ts = EXCLUDED.ts,
			level_mm = EXCLUDED.level_mm,
			density = EXCLUDED.density,
			volume_m3 = EXCLUDED.volume_m3,
			volume_l = EXCLUDED.volume_l,
			mass_t = EXCLUDED.mass_t,
			income_month = EXCLUDED.income_month,
			outflow_month = EXCLUDED.outflow_month,
			income_year = EXCLUDED.income_year,
			outflow_year = EXCLUDED.outflow_year
		WHERE boiler_readings.ts < EXCLUDED.ts
	`)
	if err != nil {
		return err
	}
	defer snapStmt.Close()

	trendStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO boiler_trend (boiler_id, ts, level_mm, volume_m3, mass_t)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (boiler_id, ts) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer trendStmt.Close()

	for _, reading := range req.Readings {
		boilerID, ok := boilerIDs[reading.BoilerNumber]
		if !ok {
			return fmt.Errorf("missing boiler mapping for number %d", reading.BoilerNumber)
		}

		if _, err := snapStmt.ExecContext(
			ctx,
			boilerID,
			reading.TS,
			reading.LevelMM,
			reading.Density,
			reading.VolumeM3,
			reading.VolumeL,
			reading.MassT,
			reading.IncomeMonth,
			reading.OutflowMonth,
			reading.IncomeYear,
			reading.OutflowYear,
		); err != nil {
			return err
		}

		if _, err := trendStmt.ExecContext(
			ctx,
			boilerID,
			reading.TS,
			reading.LevelMM,
			reading.VolumeM3,
			reading.MassT,
		); err != nil {
			return err
		}
	}

	acceptedAt := req.SentAt
	if acceptedAt.IsZero() {
		acceptedAt = time.Now().UTC()
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO plant_ingest_status (plant_id, last_accepted_at, last_count)
		VALUES ($1, $2, $3)
		ON CONFLICT (plant_id) DO UPDATE SET
			last_accepted_at = EXCLUDED.last_accepted_at,
			last_count = EXCLUDED.last_count
	`, req.PlantID, acceptedAt, len(req.Readings)); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *IngestRepo) GetPlantStatus(ctx context.Context, plantID string) (*model.PlantStatus, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT plant_id, last_accepted_at, last_count
		FROM plant_ingest_status
		WHERE plant_id = $1
	`, plantID)

	var status model.PlantStatus
	err := row.Scan(&status.PlantID, &status.LastAcceptedAt, &status.LastCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	status.LastAcceptedAt = status.LastAcceptedAt.UTC()
	return &status, nil
}
