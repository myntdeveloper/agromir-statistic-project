export interface Warehouse {
  id: number
  name: string
  total_volume_m3: number
}

export interface BoilerListItem {
  id: number
  number: number
  position: number
  level_mm: number
  volume_m3: number
  is_empty: boolean
}

export interface BoilerDetail {
  id: number
  number: number
  position: number
  warehouse_id: number
  warehouse_name: string
  diameter_mm: number
  height_max_mm: number
  updated_at: string | null
  level_mm: number
  density: number
  volume_m3: number
  volume_l: number
  mass_t: number
  income_month: number
  outflow_month: number
  income_year: number
  outflow_year: number
  is_empty: boolean
}

export interface TrendPoint {
  ts: string
  level_mm: number
  volume_m3: number
  mass_t: number
}

export interface BoilerTrend {
  boiler_id: number
  points: TrendPoint[]
}

export type Theme = 'light' | 'dark'

export interface Crumb {
  label: string
  to?: string
}
