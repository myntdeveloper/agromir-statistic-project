import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import type { Crumb } from '../types'
import { BackButton } from './BackButton'

interface Props {
  title: string
  subtitle?: ReactNode
  crumbs?: Crumb[]
  backTo?: string
  backLabel?: string
}

export function PageHeader({ title, subtitle, crumbs, backTo, backLabel }: Props) {
  const hasNav = Boolean(backTo) || Boolean(crumbs?.length)
  return (
    <header className="page-header">
      {hasNav && (
        <div className="page-header__nav">
          {backTo && <BackButton to={backTo} label={backLabel} />}
          {crumbs && crumbs.length > 0 && (
            <nav className="breadcrumbs" aria-label="Навигация">
              {crumbs.map((c, i) => (
                <span key={`${c.label}-${i}`} className="breadcrumbs__item">
                  {i > 0 && <span className="breadcrumbs__sep" aria-hidden="true">/</span>}
                  {c.to ? <Link to={c.to}>{c.label}</Link> : <span aria-current="page">{c.label}</span>}
                </span>
              ))}
            </nav>
          )}
        </div>
      )}
      <h1 className="page-header__title">{title}</h1>
      {subtitle && <div className="page-header__subtitle">{subtitle}</div>}
    </header>
  )
}
