import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { useTheme } from '../theme/useTheme'
import type { Theme, TrendPoint } from '../types'
import { formatDateTime, formatHour, formatNumber } from '../utils/format'

type SeriesKey = 'level_mm' | 'volume_m3' | 'mass_t'

interface Series {
  key: SeriesKey
  name: string
  unit: string
  digits: number
}

const SERIES: Series[] = [
  { key: 'level_mm', name: 'Уровень', unit: 'мм', digits: 0 },
  { key: 'volume_m3', name: 'Объём', unit: 'м³', digits: 2 },
  { key: 'mass_t', name: 'Масса', unit: 'т', digits: 2 },
]

const PALETTE: Record<Theme, Record<SeriesKey, string> & { grid: string; axis: string; tooltipBg: string; tooltipBorder: string; text: string }> = {
  light: {
    level_mm: '#1e5bb8',
    volume_m3: '#11866f',
    mass_t: '#b86e0b',
    grid: '#d3dae3',
    axis: '#5a6574',
    tooltipBg: '#ffffff',
    tooltipBorder: '#d3dae3',
    text: '#17202b',
  },
  dark: {
    level_mm: '#5b9bf0',
    volume_m3: '#3cc0a8',
    mass_t: '#e8a23c',
    grid: '#2f3946',
    axis: '#9aa6b6',
    tooltipBg: '#1b2129',
    tooltipBorder: '#3f4b5b',
    text: '#e4e9f0',
  },
}

const TICK_INTERVALS = 4

/** Y range snapped to a "round" step (1/2/2.5/5 × 10ⁿ) so tick labels read cleanly. */
function niceScale(values: number[], digits: number): { domain: [number, number]; ticks: number[] } {
  let min = Math.min(...values)
  let max = Math.max(...values)
  if (max - min < 1e-9) {
    const pad = Math.abs(max) * 0.01 || 1
    min -= pad
    max += pad
  }
  // a step finer than the displayed precision would produce duplicate tick labels
  const rawStep = Math.max((max - min) / TICK_INTERVALS, 10 ** -digits)
  const magnitude = 10 ** Math.floor(Math.log10(rawStep))
  const step = [1, 2, 2.5, 5, 10].map((m) => m * magnitude).find((s) => s >= rawStep) ?? 10 * magnitude

  const lo = Math.max(0, Math.floor(min / step) * step)
  const hi = Math.ceil(max / step) * step
  const ticks: number[] = []
  for (let v = lo; v <= hi + step / 2; v += step) ticks.push(Number(v.toFixed(6)))
  return { domain: [lo, hi], ticks }
}

interface Props {
  points: TrendPoint[]
}

export function BoilerTrendChart({ points }: Props) {
  const { theme } = useTheme()
  const colors = PALETTE[theme]

  if (points.length === 0) {
    return <div className="status-view">Нет данных тренда за последние сутки</div>
  }

  const axisTick = { fill: colors.axis, fontSize: 12 }

  return (
    <div className="trend">
      {SERIES.map((s) => {
        const values = points.map((p) => p[s.key])
        const current = values[values.length - 1]
        const scale = niceScale(values, s.digits)
        const fmt = (v: number) => formatNumber(v, s.digits)

        return (
          <figure key={s.key} className="trend-chart">
            <figcaption className="trend-chart__head">
              <span className="trend-chart__title">
                <span className="trend-chart__dot" style={{ background: colors[s.key] }} />
                {s.name}, {s.unit}
              </span>
              <span className="trend-chart__stats num">
                сейчас <strong>{fmt(current)}</strong>
                <span className="trend-chart__sep">·</span>
                мин {fmt(Math.min(...values))}
                <span className="trend-chart__sep">·</span>
                макс {fmt(Math.max(...values))}
              </span>
            </figcaption>

            <div className="trend-chart__body">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={points} syncId="boiler-trend" margin={{ top: 4, right: 16, bottom: 20, left: 8 }}>
                  <CartesianGrid stroke={colors.grid} strokeDasharray="3 3" vertical={false} />
                  <XAxis
                    dataKey="ts"
                    tickFormatter={formatHour}
                    tick={axisTick}
                    stroke={colors.grid}
                    minTickGap={24}
                    label={{ value: 'Время', position: 'insideBottom', offset: -14, fill: colors.axis, fontSize: 12 }}
                  />
                  <YAxis
                    domain={scale.domain}
                    ticks={scale.ticks}
                    tickFormatter={fmt}
                    tick={axisTick}
                    stroke={colors.grid}
                    width={76}
                    label={{
                      value: s.unit,
                      angle: -90,
                      position: 'insideLeft',
                      fill: colors.axis,
                      fontSize: 12,
                      style: { textAnchor: 'middle' },
                    }}
                  />
                  <Tooltip
                    contentStyle={{
                      background: colors.tooltipBg,
                      border: `1px solid ${colors.tooltipBorder}`,
                      borderRadius: 8,
                      color: colors.text,
                    }}
                    labelStyle={{ color: colors.axis }}
                    labelFormatter={(label) => formatDateTime(String(label))}
                    formatter={(value) => [`${fmt(Number(value))} ${s.unit}`, s.name]}
                  />
                  <Legend verticalAlign="top" align="right" height={24} iconType="plainline" wrapperStyle={{ fontSize: 13, color: colors.axis }} />
                  <Line
                    type="monotone"
                    dataKey={s.key}
                    name={`${s.name}, ${s.unit}`}
                    stroke={colors[s.key]}
                    strokeWidth={2.25}
                    dot={false}
                    activeDot={{ r: 4 }}
                    isAnimationActive={false}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </figure>
        )
      })}
    </div>
  )
}
