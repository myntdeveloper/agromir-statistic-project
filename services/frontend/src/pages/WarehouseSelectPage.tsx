import { api } from '../api/client'
import { PageHeader } from '../components/PageHeader'
import { ErrorView, LoadingView } from '../components/StatusView'
import { WarehouseCard } from '../components/WarehouseCard'
import { useApi } from '../hooks/useApi'

export function WarehouseSelectPage() {
  const { data, error, loading, reload } = useApi('warehouses', (signal) => api.getWarehouses(signal), {
    refreshMs: 60_000,
  })

  return (
    <div className="page">
      <PageHeader title="Выбор склада" subtitle="Ёмкости сырья по складам" />

      {loading && <LoadingView />}
      {error && <ErrorView error={error} onRetry={reload} />}
      {data && (
        <div className="warehouse-grid">
          {data.map((w) => (
            <WarehouseCard key={w.id} warehouse={w} />
          ))}
        </div>
      )}
    </div>
  )
}
