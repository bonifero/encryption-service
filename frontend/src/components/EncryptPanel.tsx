import { useState } from 'react'
import { ApiError, encryptMessage } from '../api/client'
import type { DecryptRequest, EncryptResponse } from '../api/types'
import { CopyField } from './CopyField'

interface EncryptPanelProps {
  onUseInDecrypt: (payload: DecryptRequest) => void
}

export function EncryptPanel({ onUseInDecrypt }: EncryptPanelProps) {
  const [message, setMessage] = useState('')
  const [result, setResult] = useState<EncryptResponse | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const response = await encryptMessage(message)
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
      <h2>Шифрование сообщения</h2>
      <p className="panel__hint">AES-256/GCM + RSA-OAEP(SHA-256) — гибридное шифрование на стороне сервиса.</p>
      <form onSubmit={handleSubmit} className="form">
        <label className="form__label" htmlFor="encrypt-message">
          Сообщение
        </label>
        <textarea
          id="encrypt-message"
          className="form__textarea"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Введите текст для шифрования"
          rows={4}
          required
        />
        <button type="submit" className="button" disabled={loading || message.trim() === ''}>
          {loading ? 'Шифруем…' : 'Зашифровать'}
        </button>
      </form>

      {error && <p className="error">{error}</p>}

      {result && (
        <div className="result">
          <CopyField label="encryptedKey" value={result.encryptedKey} />
          <CopyField label="iv" value={result.iv} />
          <CopyField label="cipherText" value={result.cipherText} />
          <button type="button" className="button button--secondary" onClick={() => onUseInDecrypt(result)}>
            Отправить в форму дешифрования →
          </button>
        </div>
      )}
    </section>
  )
}
