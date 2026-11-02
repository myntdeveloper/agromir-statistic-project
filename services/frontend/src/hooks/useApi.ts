import { useCallback, useEffect, useLayoutEffect, useRef, useState } from 'react'

export interface ApiState<T> {
  data: T | undefined
  error: Error | undefined
  loading: boolean
  reload: () => void
}

interface Options {
  /** Silent background refresh interval; keeps current data on screen while refetching. */
  refreshMs?: number
}

interface Result<T> {
  key: string
  data?: T
  error?: Error
}

/**
 * Loads data for `key`; whenever `key` changes the previous result is discarded
 * and the loader runs again.
 */
export function useApi<T>(
  key: string,
  loader: (signal: AbortSignal) => Promise<T>,
  { refreshMs }: Options = {},
): ApiState<T> {
  const [reloadCount, setReloadCount] = useState(0)
  const [result, setResult] = useState<Result<T>>()
  const requestKey = `${key}#${reloadCount}`

  const loaderRef = useRef(loader)
  useLayoutEffect(() => {
    loaderRef.current = loader
  })

  useEffect(() => {
    const controller = new AbortController()
    loaderRef
      .current(controller.signal)
      .then((data) => {
        if (!controller.signal.aborted) setResult({ key: requestKey, data })
      })
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setResult({ key: requestKey, error: err instanceof Error ? err : new Error(String(err)) })
      })
    return () => controller.abort()
  }, [requestKey])

  useEffect(() => {
    if (!refreshMs) return
    let controller: AbortController | undefined
    const timer = window.setInterval(() => {
      controller?.abort()
      const current = new AbortController()
      controller = current
      loaderRef
        .current(current.signal)
        .then((data) => {
          if (!current.signal.aborted) setResult({ key: requestKey, data })
        })
        .catch(() => {
          // keep the last good data on background refresh failures
        })
    }, refreshMs)

    return () => {
      window.clearInterval(timer)
      controller?.abort()
    }
  }, [requestKey, refreshMs])

  const reload = useCallback(() => setReloadCount((n) => n + 1), [])

  const current = result?.key === requestKey ? result : undefined
  return {
    data: current?.data,
    error: current?.error,
    loading: current === undefined,
    reload,
  }
}
