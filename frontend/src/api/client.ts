import type {
  ApiErrorBody,
  CryptoOperationResponse,
  DecryptRequest,
  DecryptResponse,
  EncryptResponse,
  HashAlgorithm,
  HashResponse,
  OperationType,
  PageResponse,
} from './types'

const API_BASE = '/api/crypto'

export class ApiError extends Error {
  readonly status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
    this.name = 'ApiError'
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...init?.headers },
  })

  if (!response.ok) {
    let message = `HTTP ${response.status}`
    try {
      const body = (await response.json()) as ApiErrorBody
      message = body.error || message
    } catch {
    }
    throw new ApiError(response.status, message)
  }

  return (await response.json()) as T
}

export function encryptMessage(message: string): Promise<EncryptResponse> {
  return request(`${API_BASE}/encrypt`, {
    method: 'POST',
    body: JSON.stringify({ message }),
  })
}

export function decryptMessage(payload: DecryptRequest): Promise<DecryptResponse> {
  return request(`${API_BASE}/decrypt`, {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function hashMessage(message: string, algorithm: HashAlgorithm): Promise<HashResponse> {
  return request(`${API_BASE}/hash`, {
    method: 'POST',
    body: JSON.stringify({ message, algorithm }),
  })
}

export interface LogsQuery {
  type?: OperationType
  page?: number
  size?: number
}

export function fetchLogs(query: LogsQuery): Promise<PageResponse> {
  const params = new URLSearchParams()
  if (query.type) params.set('type', query.type)
  params.set('page', String(query.page ?? 0))
  params.set('size', String(query.size ?? 10))

  return request(`${API_BASE}/logs?${params.toString()}`)
}

export function fetchLogById(id: number): Promise<CryptoOperationResponse> {
  return request(`${API_BASE}/logs/${id}`)
}
