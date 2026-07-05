package openapi

var Spec = []byte(`{
  "openapi": "3.0.3",
  "info": {
    "title": "Crypto Service API (Go)",
    "description": "Микросервис шифрования/дешифрования сообщений и расчета хеша (тестовое задание, реализация на Go)",
    "version": "1.0.0"
  },
  "servers": [{"url": "http://localhost:8081"}],
  "tags": [{"name": "Crypto", "description": "Шифрование, дешифрование и расчет хеша сообщений"}],
  "paths": {
    "/api/crypto/encrypt": {
      "post": {
        "tags": ["Crypto"],
        "summary": "Зашифровать сообщение",
        "description": "Гибридное шифрование AES-256/GCM + RSA-OAEP. Возвращает зашифрованный AES-ключ, IV и шифротекст.",
        "requestBody": {
          "required": true,
          "content": {"application/json": {"schema": {"$ref": "#/components/schemas/MessageRequest"}}}
        },
        "responses": {
          "200": {"description": "OK", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/EncryptResponse"}}}},
          "400": {"description": "Некорректный запрос", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}}}}
        }
      }
    },
    "/api/crypto/decrypt": {
      "post": {
        "tags": ["Crypto"],
        "summary": "Расшифровать сообщение",
        "description": "Дешифрует сообщение, зашифрованное методом /encrypt приватным RSA-ключом сервиса.",
        "requestBody": {
          "required": true,
          "content": {"application/json": {"schema": {"$ref": "#/components/schemas/DecryptRequest"}}}
        },
        "responses": {
          "200": {"description": "OK", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/DecryptResponse"}}}},
          "400": {"description": "Некорректный запрос", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}}}}
        }
      }
    },
    "/api/crypto/hash": {
      "post": {
        "tags": ["Crypto"],
        "summary": "Рассчитать хеш сообщения",
        "description": "Поддерживаемые алгоритмы: SHA_256, PBKDF2 (со случайной солью).",
        "requestBody": {
          "required": true,
          "content": {"application/json": {"schema": {"$ref": "#/components/schemas/HashRequest"}}}
        },
        "responses": {
          "200": {"description": "OK", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/HashResponse"}}}},
          "400": {"description": "Некорректный запрос", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}}}}
        }
      }
    },
    "/api/crypto/logs": {
      "get": {
        "tags": ["Crypto"],
        "summary": "Получить журнал выполненных криптографических операций",
        "description": "Постраничный список записей с входными данными и результатом каждой операции. Параметр type опционален (ENCRYPT, DECRYPT, HASH).",
        "parameters": [
          {"name": "type", "in": "query", "required": false, "schema": {"type": "string", "enum": ["ENCRYPT", "DECRYPT", "HASH"]}},
          {"name": "page", "in": "query", "required": false, "schema": {"type": "integer", "default": 0}},
          {"name": "size", "in": "query", "required": false, "schema": {"type": "integer", "default": 20}}
        ],
        "responses": {
          "200": {"description": "OK", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/PageResponse"}}}}
        }
      }
    },
    "/api/crypto/logs/{id}": {
      "get": {
        "tags": ["Crypto"],
        "summary": "Получить одну запись журнала операций по идентификатору",
        "parameters": [
          {"name": "id", "in": "path", "required": true, "schema": {"type": "integer"}}
        ],
        "responses": {
          "200": {"description": "OK", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/CryptoOperationResponse"}}}},
          "404": {"description": "Запись не найдена", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}}}}
        }
      }
    }
  },
  "components": {
    "schemas": {
      "MessageRequest": {
        "type": "object",
        "required": ["message"],
        "properties": {"message": {"type": "string", "example": "Hello, world!"}}
      },
      "EncryptResponse": {
        "type": "object",
        "properties": {
          "encryptedKey": {"type": "string", "description": "AES-ключ, зашифрованный открытым RSA-ключом (Base64)"},
          "iv": {"type": "string", "description": "Вектор инициализации AES/GCM (Base64)"},
          "cipherText": {"type": "string", "description": "Зашифрованное сообщение (Base64)"}
        }
      },
      "DecryptRequest": {
        "type": "object",
        "required": ["encryptedKey", "iv", "cipherText"],
        "properties": {
          "encryptedKey": {"type": "string"},
          "iv": {"type": "string"},
          "cipherText": {"type": "string"}
        }
      },
      "DecryptResponse": {
        "type": "object",
        "properties": {"message": {"type": "string"}}
      },
      "HashRequest": {
        "type": "object",
        "required": ["message", "algorithm"],
        "properties": {
          "message": {"type": "string", "example": "Hello, world!"},
          "algorithm": {"type": "string", "enum": ["SHA_256", "PBKDF2"]}
        }
      },
      "HashResponse": {
        "type": "object",
        "properties": {
          "algorithm": {"type": "string", "example": "SHA-256"},
          "hashHex": {"type": "string"},
          "saltHex": {"type": "string", "nullable": true},
          "iterations": {"type": "integer", "nullable": true}
        }
      },
      "CryptoOperationResponse": {
        "type": "object",
        "properties": {
          "id": {"type": "integer"},
          "operationType": {"type": "string", "enum": ["ENCRYPT", "DECRYPT", "HASH"]},
          "inputData": {"type": "string"},
          "outputData": {"type": "string"},
          "createdAt": {"type": "string", "format": "date-time"}
        }
      },
      "PageResponse": {
        "type": "object",
        "properties": {
          "content": {"type": "array", "items": {"$ref": "#/components/schemas/CryptoOperationResponse"}},
          "totalElements": {"type": "integer"},
          "totalPages": {"type": "integer"},
          "size": {"type": "integer"},
          "number": {"type": "integer"},
          "numberOfElements": {"type": "integer"},
          "first": {"type": "boolean"},
          "last": {"type": "boolean"},
          "empty": {"type": "boolean"}
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "timestamp": {"type": "string"},
          "status": {"type": "integer"},
          "error": {"type": "string"}
        }
      }
    }
  }
}`)
