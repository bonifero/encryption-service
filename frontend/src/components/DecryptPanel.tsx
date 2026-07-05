import { useEffect, useState } from 'react'
import { ApiError, decryptMessage } from '../api/client'
import type { DecryptRequest } from '../api/types'

interface DecryptPanelProps {
  prefill: DecryptRequest | null
}

const EMPTY: DecryptRequest = { encryptedKey: '', iv: '', cipherText: '' }

export function DecryptPanel({ prefill }: DecryptPanelProps) {
  const [form, setForm] = useState<DecryptRequest>(EMPTY)
  const [result, setResult] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (prefill) {
      setForm(prefill)
      setResult(null)
      setError(null)
    }
  }, [prefill])

  function update(field: keyof DecryptRequest, value: string) {
    setForm((prev) => ({ ...prev, [field]: value }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const response = await decryptMessage(form)
      setResult(response.message)
    } catch (err) {
      setResult(null)
      setError(err instanceof ApiError ? err.message : 'Не удалось выполнить запрос')
    } finally {
      setLoading(false)
    }
  }

  const isComplete = form.encryptedKey.trim() && form.iv.trim() && form.cipherText.trim()

  return (
    <section className="panel">
      <h2>Дешифрование сообщения</h2>
      <p className="panel__hint">Вставьте данные, полученные от /encrypt (приватным RSA-ключом сервиса).</p>
      <form onSubmit={handleSubmit} className="form">
        <label className="form__label" htmlFor="decrypt-key">
          encryptedKey
        </label>
        <textarea
          id="decrypt-key"
          className="form__textarea"
          value={form.encryptedKey}
          onChange={(e) => update('encryptedKey', e.target.value)}
          rows={2}
          required
        />

        <label className="form__label" htmlFor="decrypt-iv">
          iv
        </label>
        <textarea
          id="decrypt-iv"
          className="form__textarea"
          value={form.iv}
          onChange={(e) => update('iv', e.target.value)}
          rows={1}
          required
        />

        <label className="form__label" htmlFor="decrypt-ciphertext">
          cipherText
        </label>
        <textarea
          id="decrypt-ciphertext"
          className="form__textarea"
          value={form.cipherText}
          onChange={(e) => update('cipherText', e.target.value)}
          rows={2}
          required
        />

        <button type="submit" className="button" disabled={loading || !isComplete}>
          {loading ? 'Дешифруем…' : 'Расшифровать'}
        </button>
      </form>

      {error && <p className="error">{error}</p>}

      {result !== null && (
        <div className="result">
          <div className="result__box">{result}</div>
        </div>
      )}
    </section>
  )
}
