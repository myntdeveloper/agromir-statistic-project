import { ApiError } from '../api/client'

export function LoadingView({ label = 'Загрузка…' }: { label?: string }) {
  return (
    <div className="status-view" role="status">
      <span className="spinner" aria-hidden="true" />
      <span>{label}</span>
    </div>
  )
}

interface ErrorProps {
  error: Error
  onRetry?: () => void
  notFoundText?: string
}

export function ErrorView({ error, onRetry, notFoundText = 'Не найдено' }: ErrorProps) {
  const isNotFound = error instanceof ApiError && error.status === 404
  return (
    <div className="status-view status-view--error" role="alert">
      <strong>{isNotFound ? notFoundText : 'Не удалось загрузить данные'}</strong>
      {!isNotFound && <span className="status-view__details">{error.message}</span>}
      {onRetry && !isNotFound && (
        <button type="button" className="btn" onClick={onRetry}>
          Повторить
        </button>
      )}
    </div>
  )
}
