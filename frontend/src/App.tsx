import { useState } from 'react'
import './App.css'
import { DecryptPanel } from './components/DecryptPanel'
import { EncryptPanel } from './components/EncryptPanel'
import { HashPanel } from './components/HashPanel'
import { LogsPanel } from './components/LogsPanel'
import type { DecryptRequest } from './api/types'

type Tab = 'encrypt' | 'decrypt' | 'hash' | 'logs'

const TABS: { id: Tab; label: string }[] = [
  { id: 'encrypt', label: 'Шифрование' },
  { id: 'decrypt', label: 'Дешифрование' },
  { id: 'hash', label: 'Хеш' },
  { id: 'logs', label: 'Журнал' },
]

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('encrypt')
  const [decryptPrefill, setDecryptPrefill] = useState<DecryptRequest | null>(null)

  function handleUseInDecrypt(payload: DecryptRequest) {
    setDecryptPrefill(payload)
    setActiveTab('decrypt')
  }

  return (
    <div className="app">
      <header className="app__header">
        <h1>Crypto Service</h1>
        <p>Веб-клиент для сервиса шифрования, дешифрования и расчета хеша сообщений</p>
      </header>

      <nav className="tabs">
        {TABS.map((tab) => (
          <button
            key={tab.id}
            type="button"
            className={`tabs__item ${activeTab === tab.id ? 'tabs__item--active' : ''}`}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
          </button>
        ))}
      </nav>

      <main className="app__content">
        {activeTab === 'encrypt' && <EncryptPanel onUseInDecrypt={handleUseInDecrypt} />}
        {activeTab === 'decrypt' && <DecryptPanel prefill={decryptPrefill} />}
        {activeTab === 'hash' && <HashPanel />}
        {activeTab === 'logs' && <LogsPanel />}
      </main>
    </div>
  )
}

export default App
