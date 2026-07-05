import { useState } from 'react'

interface CopyFieldProps {
  label: string
  value: string
}

export function CopyField({ label, value }: CopyFieldProps) {
  const [copied, setCopied] = useState(false)

  async function handleCopy() {
    await navigator.clipboard.writeText(value)
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  return (
    <div className="copy-field">
      <div className="copy-field__label-row">
        <span className="copy-field__label">{label}</span>
        <button type="button" className="link-button" onClick={handleCopy}>
          {copied ? 'Скопировано' : 'Копировать'}
        </button>
      </div>
      <textarea className="copy-field__value" value={value} readOnly rows={2} />
    </div>
  )
}
