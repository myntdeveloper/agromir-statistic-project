INSERT INTO warehouses (id, name, sort_order) VALUES
  (1, 'Склад 1', 1),
  (2, 'Склад 2', 2),
  (3, 'Склад 3', 3);

SELECT setval(pg_get_serial_sequence('warehouses', 'id'), (SELECT MAX(id) FROM warehouses));

-- boilers: 1-8 -> склад 1, 9-16 -> склад 2, 17-24 -> склад 3
INSERT INTO boilers (id, warehouse_id, number, position)
SELECT n, 1, n, n      FROM generate_series(1, 8)   AS n
UNION ALL
SELECT n, 2, n, n - 8  FROM generate_series(9, 16)  AS n
UNION ALL
SELECT n, 3, n, n - 16 FROM generate_series(17, 24) AS n;

SELECT setval(pg_get_serial_sequence('boilers', 'id'), (SELECT MAX(id) FROM boilers));

INSERT INTO boiler_readings (boiler_id, ts, level_mm, density, volume_m3, volume_l, mass_t, income_month, outflow_month, income_year, outflow_year) VALUES
-- Склад 1: реальные тестовые показания
(1, now(), 7298, 1.35, 159.939172, 159939.17, 215.917882, 1.35, 1.35, 1.35, 1.35),
(2, now(), 7313, 1.35, 159.608242, 159608.24, 215.471127, 1.35, 1.35, 1.35, 1.35),
(3, now(), 7267, 1.35, 158.372770, 158372.77, 213.803240, 1.35, 1.35, 1.35, 1.35),
(4, now(), 2976, 1.35,  63.704728,  63704.73,  86.001383, 1.35, 1.35, 1.35, 1.35),
(5, now(), 7328, 1.35, 160.159792, 160159.79, 216.215719, 1.35, 1.35, 1.35, 1.35),
(6, now(), 7283, 1.35, 158.505142, 158505.14, 213.981942, 1.35, 1.35, 1.35, 1.35),
(7, now(), 7322, 1.35, 160.159792, 160159.79, 216.215720, 1.35, 1.35, 1.35, 1.35),
(8, now(), 7313, 1.35, 159.387622, 159387.62, 215.173290, 1.35, 1.35, 1.35, 1.35),
-- Склад 2: частично заполнен
(9,  now(), 7313, 1.35, 159.078754, 159078.75, 214.756318, 1.35, 1.35, 1.35, 1.35),
(10, now(), 7329, 1.35, 159.740614, 159740.61, 215.649829, 1.35, 1.35, 1.35, 1.35),
(11, now(),    0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(12, now(),    0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(13, now(),    0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(14, now(),    0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(15, now(),    0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(16, now(),    0, 1.35, 0, 0, 0, 0, 0, 0, 0),
-- Склад 3: пустой
(17, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(18, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(19, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(20, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(21, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(22, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(23, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0),
(24, now(), 0, 1.35, 0, 0, 0, 0, 0, 0, 0);

-- 24 hourly points per boiler: ±0.5% noise on level/volume, mass = volume * density.
-- The newest point (h = 0) matches the current reading exactly.
WITH noise AS MATERIALIZED (
    SELECT br.boiler_id,
           date_trunc('hour', now()) - make_interval(hours => h) AS ts,
           br.level_mm,
           br.volume_m3,
           br.density,
           CASE WHEN h = 0 THEN 1::numeric
                ELSE (1 + (random() - 0.5) * 0.01)::numeric
           END AS k
    FROM boiler_readings br
    CROSS JOIN generate_series(0, 23) AS h
)
INSERT INTO boiler_trend (boiler_id, ts, level_mm, volume_m3, mass_t)
SELECT boiler_id,
       ts,
       ROUND(level_mm * k),
       ROUND(volume_m3 * k, 3),
       ROUND(volume_m3 * k * density, 3)
FROM noise
ON CONFLICT (boiler_id, ts) DO NOTHING;
