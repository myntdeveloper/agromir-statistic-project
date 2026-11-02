import { Link } from 'react-router-dom'
import type { Warehouse } from '../types'
import { formatNumber } from '../utils/format'

export function WarehouseCard({ warehouse }: { warehouse: Warehouse }) {
  const isEmpty = warehouse.total_volume_m3 <= 0
  return (
    <Link to={`/warehouse/${warehouse.id}`} className={`warehouse-card${isEmpty ? ' is-empty' : ''}`}>
      <span className="warehouse-card__name">{warehouse.name}</span>
      <span className="warehouse-card__volume">
        <span className="num">{formatNumber(warehouse.total_volume_m3, 2)}</span>
        <span className="unit">м³</span>
      </span>
      <span className="warehouse-card__caption">{isEmpty ? 'пусто' : 'суммарный объём'}</span>
    </Link>
  )
}
