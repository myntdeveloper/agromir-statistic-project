import { Link } from 'react-router-dom'

interface Props {
  to: string
  label?: string
}

export function BackButton({ to, label = 'Назад' }: Props) {
  return (
    <Link to={to} className="back-button">
      <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <path d="M15 18l-6-6 6-6" />
      </svg>
      <span>{label}</span>
    </Link>
  )
}
