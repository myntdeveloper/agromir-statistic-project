import { Link } from 'react-router-dom'
import type { BoilerListItem } from '../types'
import { formatNumber } from '../utils/format'

interface Props {
  boiler: BoilerListItem
  to: string
}

export function BoilerButton({ boiler, to }: Props) {
  return (
    <Link
      to={to}
      className={`boiler-button${boiler.is_empty ? ' is-empty' : ''}`}
      aria-label={`Ёмкость ${boiler.number}${boiler.is_empty ? ', пусто' : ''}`}
    >
      <span className="boiler-button__number">{boiler.number}</span>
      <span className="boiler-button__meta">
        {boiler.is_empty ? (
          'пусто'
        ) : (
          <>
            <span className="num">{formatNumber(boiler.volume_m3, 2)}</span> м³
          </>
        )}
      </span>
    </Link>
  )
}
