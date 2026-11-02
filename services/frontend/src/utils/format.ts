const numberFormats = new Map<number, Intl.NumberFormat>()

function numberFormat(digits: number): Intl.NumberFormat {
  let nf = numberFormats.get(digits)
  if (!nf) {
    nf = new Intl.NumberFormat('en-US', {
      minimumFractionDigits: digits,
      maximumFractionDigits: digits,
    })
    numberFormats.set(digits, nf)
  }
  return nf
}

/** 159939.17 -> "159 939.17" (non-breaking space as thousands separator, dot as decimal). */
export function formatNumber(value: number, digits = 2): string {
  if (!Number.isFinite(value)) return '—'
  return numberFormat(digits).format(value).replace(/,/g, '\u00a0')
}

const dateTimeFormat = new Intl.DateTimeFormat('ru-RU', {
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
})

const hourFormat = new Intl.DateTimeFormat('ru-RU', { hour: '2-digit', minute: '2-digit' })

export function formatDateTime(iso: string | null | undefined): string {
  if (!iso) return '—'
  const d = new Date(iso)
  return Number.isNaN(d.getTime()) ? '—' : dateTimeFormat.format(d)
}

export function formatHour(iso: string): string {
  const d = new Date(iso)
  return Number.isNaN(d.getTime()) ? '' : hourFormat.format(d)
}
