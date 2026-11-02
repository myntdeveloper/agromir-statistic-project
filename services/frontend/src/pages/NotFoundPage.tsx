import { PageHeader } from '../components/PageHeader'

export function NotFoundPage() {
  return (
    <div className="page">
      <PageHeader title="Страница не найдена" backTo="/" backLabel="К складам" />
    </div>
  )
}
