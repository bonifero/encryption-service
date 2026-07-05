import { Fragment, useEffect, useState } from 'react'
import { ApiError, fetchLogs } from '../api/client'
import type { CryptoOperationResponse, OperationType, PageResponse } from '../api/types'

const PAGE_SIZE = 10

function prettyJson(raw: string): string {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

export function LogsPanel() {
  const [typeFilter, setTypeFilter] = useState<OperationType | ''>('')
  const [page, setPage] = useState(0)
  const [data, setData] = useState<PageResponse | null>(null)
  const [expandedId, setExpandedId] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let cancelled = false
    setLoading(true)
    setError(null)

    fetchLogs({ type: typeFilter || undefined, page, size: PAGE_SIZE })
      .then((response) => {
        if (!cancelled) setData(response)
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof ApiError ? err.message : 'Не удалось загрузить журнал')
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [typeFilter, page])

  function handleFilterChange(value: string) {
    setTypeFilter(value as OperationType | '')
    setPage(0)
    setExpandedId(null)
  }

  function toggleExpand(op: CryptoOperationResponse) {
    setExpandedId((current) => (current === op.id ? null : op.id))
  }

  return (
    <section className="panel">
      <h2>Журнал операций</h2>
      <p className="panel__hint">Входные данные и результат каждого вызова /encrypt, /decrypt и /hash.</p>

      <div className="logs__toolbar">
        <label htmlFor="logs-type-filter">Тип операции:</label>
        <select
          id="logs-type-filter"
          className="form__select"
          value={typeFilter}
          onChange={(e) => handleFilterChange(e.target.value)}
        >
          <option value="">Все</option>
          <option value="ENCRYPT">ENCRYPT</option>
          <option value="DECRYPT">DECRYPT</option>
          <option value="HASH">HASH</option>
        </select>
      </div>

      {error && <p className="error">{error}</p>}
      {loading && <p className="panel__hint">Загрузка…</p>}

      {data && (
        <>
          <table className="logs-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Тип</th>
                <th>Дата</th>
                <th />
              </tr>
            </thead>
            <tbody>
              {data.content.map((op) => (
                <Fragment key={op.id}>
                  <tr className="logs-table__row" onClick={() => toggleExpand(op)}>
                    <td>{op.id}</td>
                    <td>
                      <span className={`badge badge--${op.operationType.toLowerCase()}`}>{op.operationType}</span>
                    </td>
                    <td>{new Date(op.createdAt).toLocaleString()}</td>
                    <td className="logs-table__toggle">{expandedId === op.id ? '▲' : '▼'}</td>
                  </tr>
                  {expandedId === op.id && (
                    <tr>
                      <td colSpan={4}>
                        <div className="logs-detail">
                          <div>
                            <h4>Вход</h4>
                            <pre>{prettyJson(op.inputData)}</pre>
                          </div>
                          <div>
                            <h4>Результат</h4>
                            <pre>{prettyJson(op.outputData)}</pre>
                          </div>
                        </div>
                      </td>
                    </tr>
                  )}
                </Fragment>
              ))}
              {data.content.length === 0 && (
                <tr>
                  <td colSpan={4} className="logs-table__empty">
                    Записей нет
                  </td>
                </tr>
              )}
            </tbody>
          </table>

          <div className="logs__pagination">
            <button
              type="button"
              className="button button--secondary"
              disabled={data.first}
              onClick={() => setPage((p) => Math.max(0, p - 1))}
            >
              ← Назад
            </button>
            <span>
              Страница {data.number + 1} из {Math.max(data.totalPages, 1)} ({data.totalElements} записей)
            </span>
            <button
              type="button"
              className="button button--secondary"
              disabled={data.last}
              onClick={() => setPage((p) => p + 1)}
            >
              Вперёд →
            </button>
          </div>
        </>
      )}
    </section>
  )
}
