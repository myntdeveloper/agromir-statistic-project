import { useParams } from 'react-router-dom'
import { api } from '../api/client'
import { BoilerButton } from '../components/BoilerButton'
import { PageHeader } from '../components/PageHeader'
import { ErrorView, LoadingView } from '../components/StatusView'
import { useApi } from '../hooks/useApi'
import { formatNumber } from '../utils/format'

export function BoilerGridPage() {
  const warehouseId = Number(useParams().warehouseId)

  const { data, error, loading, reload } = useApi(
    `warehouse-${warehouseId}`,
    async (signal) => {
      const [warehouses, boilers] = await Promise.all([
        api.getWarehouses(signal),
        api.getWarehouseBoilers(warehouseId, signal),
      ])
      return { warehouse: warehouses.find((w) => w.id === warehouseId), boilers }
    },
    { refreshMs: 60_000 },
  )

  const title = data?.warehouse?.name ?? `Склад ${warehouseId}`
  const filled = data?.boilers.filter((b) => !b.is_empty).length ?? 0

  return (
    <div className="page">
      <PageHeader
        title={title}
        backTo="/"
        backLabel="К складам"
        crumbs={[{ label: 'Склады', to: '/' }, { label: title }]}
        subtitle={
          data?.warehouse && (
            <>
              Суммарный объём <span className="num">{formatNumber(data.warehouse.total_volume_m3, 2)}</span> м³ · заполнено{' '}
              {filled} из {data.boilers.length}
            </>
          )
        }
      />

      {loading && <LoadingView />}
      {error && <ErrorView error={error} onRetry={reload} notFoundText="Склад не найден" />}
      {data && (
        <div className="boiler-grid">
          {data.boilers.map((b) => (
            <BoilerButton key={b.id} boiler={b} to={`/warehouse/${warehouseId}/boiler/${b.id}`} />
          ))}
        </div>
      )}
    </div>
  )
}
