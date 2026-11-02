import type { BoilerDetail } from '../types'
import { formatDateTime, formatNumber } from '../utils/format'

interface Row {
  label: string
  value: number
}

export function BoilerInfoPanel({ boiler }: { boiler: BoilerDetail }) {
  const fillPct =
    boiler.height_max_mm > 0 ? Math.min(100, Math.max(0, (boiler.level_mm / boiler.height_max_mm) * 100)) : 0

  const rows: Row[] = [
    { label: 'Объём, м³', value: boiler.volume_m3 },
    { label: 'Объём, л', value: boiler.volume_l },
    { label: 'Масса в Т', value: boiler.mass_t },
    { label: 'Приход за месяц', value: boiler.income_month },
    { label: 'Расход за месяц', value: boiler.outflow_month },
    { label: 'Приход за год', value: boiler.income_year },
    { label: 'Расход за год', value: boiler.outflow_year },
  ]

  return (
    <section className="info-panel" aria-label="Информация о ёмкости">
      <div className="info-panel__summary">
        <div className="stat stat--primary">
          <span className="stat__label">Уровень налива H</span>
          <span className="stat__value">
            <span className="num">{formatNumber(boiler.level_mm, 0)}</span>
            <span className="unit">мм</span>
          </span>
          <div
            className="fill-bar"
            role="meter"
            aria-label="Заполнение ёмкости"
            aria-valuemin={0}
            aria-valuemax={100}
            aria-valuenow={Math.round(fillPct)}
          >
            <div className="fill-bar__value" style={{ width: `${fillPct}%` }} />
          </div>
          <span className="stat__hint">
            {formatNumber(fillPct, 0)}% от {formatNumber(boiler.height_max_mm, 0)} мм
          </span>
        </div>

        <div className="stat">
          <span className="stat__label">Плотность</span>
          <span className="stat__value">
            <span className="num">{formatNumber(boiler.density, 2)}</span>
            <span className="unit">т/м³</span>
          </span>
        </div>

        <div className="stat">
          <span className="stat__label">Обновлено</span>
          <span className="stat__value stat__value--small">
            <span className="num">{formatDateTime(boiler.updated_at)}</span>
          </span>
        </div>
      </div>

      <table className="info-table">
        <tbody>
          {rows.map((row) => (
            <tr key={row.label}>
              <th scope="row">{row.label}:</th>
              <td className="num">{formatNumber(row.value, 2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}
