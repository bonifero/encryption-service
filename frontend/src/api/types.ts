export type HashAlgorithm = 'SHA_256' | 'PBKDF2'

export type OperationType = 'ENCRYPT' | 'DECRYPT' | 'HASH'

export interface EncryptResponse {
  encryptedKey: string
  iv: string
  cipherText: string
}

export interface DecryptRequest {
  encryptedKey: string
  iv: string
  cipherText: string
}

export interface DecryptResponse {
  message: string
}

export interface HashResponse {
  algorithm: string
  hashHex: string
  saltHex: string | null
  iterations: number | null
}

export interface CryptoOperationResponse {
  id: number
  operationType: OperationType
  inputData: string
  outputData: string
  createdAt: string
}

export interface PageResponse {
  content: CryptoOperationResponse[]
  totalElements: number
  totalPages: number
  size: number
  number: number
  numberOfElements: number
  first: boolean
  last: boolean
  empty: boolean
}

export interface ApiErrorBody {
  timestamp: string
  status: number
  error: string
}
