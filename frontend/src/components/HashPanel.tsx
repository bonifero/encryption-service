import { useState } from 'react'
import { ApiError, hashMessage } from '../api/client'
import type { HashAlgorithm, HashResponse } from '../api/types'

export function HashPanel() {
  const [message, setMessage] = useState('')
  const [algorithm, setAlgorithm] = useState<HashAlgorithm>('SHA_256')
  const [result, setResult] = useState<HashResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const response = await hashMessage(message, algorithm)
      setResult(response)
    } catch (err) {
      setResult(null)
      setError(err instanceof ApiError ? err.message : 'Не удалось выполнить запрос')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="panel">
      <h2>Расчет хеша сообщения</h2>
      <p className="panel__hint">SHA-256 — детерминированный. PBKDF2WithHmacSHA256 — со случайной солью (100000 итераций).</p>
      <form onSubmit={handleSubmit} className="form">
        <label className="form__label" htmlFor="hash-message">
          Сообщение
        </label>
        <textarea
          id="hash-message"
          className="form__textarea"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Введите текст для расчета хеша"
          rows={4}
          required
        />

        <label className="form__label" htmlFor="hash-algorithm">
          Алгоритм
        </label>
        <select
          id="hash-algorithm"
          className="form__select"
          value={algorithm}
          onChange={(e) => setAlgorithm(e.target.value as HashAlgorithm)}
        >
          <option value="SHA_256">SHA-256</option>
          <option value="PBKDF2">PBKDF2WithHmacSHA256</option>
        </select>

        <button type="submit" className="button" disabled={loading || message.trim() === ''}>
          {loading ? 'Считаем…' : 'Рассчитать хеш'}
        </button>
      </form>

      {error && <p className="error">{error}</p>}

      {result && (
        <div className="result">
          <dl className="kv">
            <dt>Алгоритм</dt>
            <dd>{result.algorithm}</dd>
            <dt>Хеш (hex)</dt>
            <dd className="kv__mono">{result.hashHex}</dd>
            {result.saltHex && (
              <>
                <dt>Соль (hex)</dt>
                <dd className="kv__mono">{result.saltHex}</dd>
              </>
            )}
            {result.iterations !== null && (
              <>
                <dt>Итераций</dt>
                <dd>{result.iterations}</dd>
              </>
            )}
          </dl>
        </div>
      )}
    </section>
  )
}
