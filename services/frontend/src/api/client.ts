import type { BoilerDetail, BoilerListItem, BoilerTrend, Warehouse } from '../types'

const BASE_URL = (import.meta.env.VITE_API_URL || '/api/v1').replace(/\/+$/, '')

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

async function request<T>(path: string, signal?: AbortSignal): Promise<T> {
  let res: Response
  try {
    res = await fetch(`${BASE_URL}${path}`, {
      signal,
      headers: { Accept: 'application/json' },
    })
  } catch (err) {
    if (signal?.aborted) throw err
    throw new ApiError(0, 'Сервер недоступен')
  }

  if (!res.ok) {
    let message = `Ошибка сервера (${res.status})`
    try {
      const body = (await res.json()) as { error?: string }
      if (body?.error) message = body.error
    } catch {
      // non-JSON error body
    }
    throw new ApiError(res.status, message)
  }

  return (await res.json()) as T
}

export const api = {
  getWarehouses: (signal?: AbortSignal) => request<Warehouse[]>('/warehouses', signal),

  getWarehouseBoilers: (warehouseId: number, signal?: AbortSignal) =>
    request<BoilerListItem[]>(`/warehouses/${warehouseId}/boilers`, signal),

  getBoiler: (boilerId: number, signal?: AbortSignal) =>
    request<BoilerDetail>(`/boilers/${boilerId}`, signal),

  getBoilerTrend: (boilerId: number, signal?: AbortSignal) =>
    request<BoilerTrend>(`/boilers/${boilerId}/trend`, signal),
}
