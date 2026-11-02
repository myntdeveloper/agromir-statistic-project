CREATE TABLE IF NOT EXISTS warehouses (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    sort_order  INT NOT NULL
);

CREATE TABLE IF NOT EXISTS boilers (
    id            SERIAL PRIMARY KEY,
    warehouse_id  INT NOT NULL REFERENCES warehouses(id),
    number        INT NOT NULL UNIQUE,
    position      INT NOT NULL,
    diameter_mm   INT NOT NULL DEFAULT 5300,
    height_max_mm INT NOT NULL DEFAULT 7400
);

CREATE INDEX IF NOT EXISTS idx_boilers_warehouse ON boilers (warehouse_id, position);

CREATE TABLE IF NOT EXISTS boiler_readings (
    boiler_id      INT PRIMARY KEY REFERENCES boilers(id),
    ts             TIMESTAMPTZ NOT NULL,
    level_mm       NUMERIC NOT NULL,
    density        NUMERIC NOT NULL,
    volume_m3      NUMERIC NOT NULL,
    volume_l       NUMERIC NOT NULL,
    mass_t         NUMERIC NOT NULL,
    income_month   NUMERIC NOT NULL DEFAULT 0,
    outflow_month  NUMERIC NOT NULL DEFAULT 0,
    income_year    NUMERIC NOT NULL DEFAULT 0,
    outflow_year   NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS boiler_trend (
    id         BIGSERIAL PRIMARY KEY,
    boiler_id  INT NOT NULL REFERENCES boilers(id),
    ts         TIMESTAMPTZ NOT NULL,
    level_mm   NUMERIC NOT NULL,
    volume_m3  NUMERIC NOT NULL,
    mass_t     NUMERIC NOT NULL,
    UNIQUE (boiler_id, ts)
);

CREATE INDEX IF NOT EXISTS idx_boiler_trend_boiler_ts ON boiler_trend (boiler_id, ts DESC);

CREATE TABLE IF NOT EXISTS plant_ingest_status (
    plant_id          TEXT PRIMARY KEY,
    last_accepted_at  TIMESTAMPTZ NOT NULL,
    last_count        INT NOT NULL
);
