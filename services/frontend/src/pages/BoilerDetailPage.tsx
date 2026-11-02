import { useParams, useSearchParams } from 'react-router-dom'
import { api } from '../api/client'
import { BoilerInfoPanel } from '../components/BoilerInfoPanel'
import { BoilerTrendChart } from '../components/BoilerTrendChart'
import { PageHeader } from '../components/PageHeader'
import { ErrorView, LoadingView } from '../components/StatusView'
import { useApi } from '../hooks/useApi'

type Tab = 'info' | 'chart'

const TABS: { id: Tab; label: string }[] = [
  { id: 'info', label: 'Информация' },
  { id: 'chart', label: 'График' },
]

export function BoilerDetailPage() {
  const params = useParams()
  const warehouseId = Number(params.warehouseId)
  const boilerId = Number(params.boilerId)

  const [searchParams, setSearchParams] = useSearchParams()
  const tab: Tab = searchParams.get('tab') === 'chart' ? 'chart' : 'info'

  const selectTab = (next: Tab) => {
    if (next === tab) return
    setSearchParams(next === 'chart' ? { tab: 'chart' } : {}, { replace: true })
  }

  const boiler = useApi(`boiler-${boilerId}`, (signal) => api.getBoiler(boilerId, signal), { refreshMs: 60_000 })
  const trend = useApi(`trend-${boilerId}`, (signal) => api.getBoilerTrend(boilerId, signal), {
    refreshMs: 300_000,
  })

  const warehouseName = boiler.data?.warehouse_name ?? `Склад ${warehouseId}`
  const title = boiler.data ? `Ёмкость №${boiler.data.number}` : 'Ёмкость'

  return (
    <div className="page">
      <PageHeader
        title={title}
        subtitle={boiler.data ? warehouseName : undefined}
        backTo={`/warehouse/${warehouseId}`}
        backLabel="К ёмкостям"
        crumbs={[
          { label: 'Склады', to: '/' },
          { label: warehouseName, to: `/warehouse/${warehouseId}` },
          { label: title },
        ]}
      />

      <div className="tabs" role="tablist" aria-label="Режим просмотра">
        {TABS.map((t) => (
          <button
            key={t.id}
            type="button"
            role="tab"
            id={`tab-${t.id}`}
            aria-selected={tab === t.id}
            aria-controls={`panel-${t.id}`}
            className={`tabs__tab${tab === t.id ? ' is-active' : ''}`}
            onClick={() => selectTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </div>

      {tab === 'info' && (
        <div role="tabpanel" id="panel-info" aria-labelledby="tab-info" className="tab-panel">
          {boiler.loading && <LoadingView />}
          {boiler.error && <ErrorView error={boiler.error} onRetry={boiler.reload} notFoundText="Ёмкость не найдена" />}
          {boiler.data && <BoilerInfoPanel boiler={boiler.data} />}
        </div>
      )}

      {tab === 'chart' && (
        <div role="tabpanel" id="panel-chart" aria-labelledby="tab-chart" className="tab-panel">
          <h2 className="section-title">Тренд уровня, объёма, массы за сутки</h2>
          {trend.loading && <LoadingView />}
          {trend.error && <ErrorView error={trend.error} onRetry={trend.reload} notFoundText="Ёмкость не найдена" />}
          {trend.data && <BoilerTrendChart points={trend.data.points} />}
        </div>
      )}
    </div>
  )
}
