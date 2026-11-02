package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/myntdeveloper/agromir_boiler_server/internal/model"
)

type BoilerRepo struct {
	db *sql.DB
}

func NewBoilerRepo(db *sql.DB) *BoilerRepo {
	return &BoilerRepo{db: db}
}

func (r *BoilerRepo) ListWarehouses(ctx context.Context) ([]model.WarehouseSummary, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT w.id, w.name, ROUND(COALESCE(SUM(br.volume_m3), 0), 2)
		FROM warehouses w
		LEFT JOIN boilers b ON b.warehouse_id = w.id
		LEFT JOIN boiler_readings br ON br.boiler_id = b.id
		GROUP BY w.id, w.name, w.sort_order
		ORDER BY w.sort_order, w.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.WarehouseSummary, 0)
	for rows.Next() {
		var w model.WarehouseSummary
		if err := rows.Scan(&w.ID, &w.Name, &w.TotalVolumeM3); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (r *BoilerRepo) WarehouseExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM warehouses WHERE id = $1)`, id).Scan(&exists)
	return exists, err
}

func (r *BoilerRepo) BoilerExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM boilers WHERE id = $1)`, id).Scan(&exists)
	return exists, err
}

func (r *BoilerRepo) ListBoilersByWarehouse(ctx context.Context, warehouseID int) ([]model.BoilerListItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.id, b.number, b.position,
		       COALESCE(br.level_mm, 0),
		       ROUND(COALESCE(br.volume_m3, 0), 2)
		FROM boilers b
		LEFT JOIN boiler_readings br ON br.boiler_id = b.id
		WHERE b.warehouse_id = $1
		ORDER BY b.position, b.number
	`, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.BoilerListItem, 0)
	for rows.Next() {
		var b model.BoilerListItem
		if err := rows.Scan(&b.ID, &b.Number, &b.Position, &b.LevelMM, &b.VolumeM3); err != nil {
			return nil, err
		}
		b.IsEmpty = b.LevelMM <= 0
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *BoilerRepo) GetBoiler(ctx context.Context, id int) (*model.BoilerDetail, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT b.id, b.number, b.position, w.id, w.name, b.diameter_mm, b.height_max_mm,
		       br.ts,
		       COALESCE(br.level_mm, 0), COALESCE(br.density, 0),
		       COALESCE(br.volume_m3, 0), COALESCE(br.volume_l, 0), COALESCE(br.mass_t, 0),
		       COALESCE(br.income_month, 0), COALESCE(br.outflow_month, 0),
		       COALESCE(br.income_year, 0), COALESCE(br.outflow_year, 0)
		FROM boilers b
		JOIN warehouses w ON w.id = b.warehouse_id
		LEFT JOIN boiler_readings br ON br.boiler_id = b.id
		WHERE b.id = $1
	`, id)

	var d model.BoilerDetail
	var ts sql.NullTime
	err := row.Scan(
		&d.ID, &d.Number, &d.Position, &d.WarehouseID, &d.WarehouseName, &d.DiameterMM, &d.HeightMaxMM,
		&ts,
		&d.LevelMM, &d.Density,
		&d.VolumeM3, &d.VolumeL, &d.MassT,
		&d.IncomeMonth, &d.OutflowMonth,
		&d.IncomeYear, &d.OutflowYear,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if ts.Valid {
		t := ts.Time.UTC()
		d.UpdatedAt = &t
	}
	d.IsEmpty = d.LevelMM <= 0
	return &d, nil
}

// GetTrend returns the 24-hour window ending at the boiler's most recent trend point,
// so the chart stays populated even when the data source has paused.
func (r *BoilerRepo) GetTrend(ctx context.Context, boilerID int) ([]model.TrendPoint, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ts, level_mm, volume_m3, mass_t
		FROM boiler_trend
		WHERE boiler_id = $1
		  AND ts > (SELECT MAX(ts) FROM boiler_trend WHERE boiler_id = $1) - INTERVAL '24 hours'
		ORDER BY ts ASC
	`, boilerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.TrendPoint, 0, 24)
	for rows.Next() {
		var p model.TrendPoint
		if err := rows.Scan(&p.TS, &p.LevelMM, &p.VolumeM3, &p.MassT); err != nil {
			return nil, err
		}
		p.TS = p.TS.UTC()
		out = append(out, p)
	}
	return out, rows.Err()
}
