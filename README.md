# Crypto Service — тестовое задание

Микросервис шифрования, дешифрования и расчета хеша сообщений с REST API,
веб-клиентом и журналом операций в PostgreSQL. Реализован в двух independent
вариантах бэкенда с одинаковым API-контрактом — **Java** (`crypto/`) и
**Go** (`crypto-go/`) — плюс общий фронтенд (`frontend/`) и nginx как
единая точка входа (`nginx/`).

Бэкенды взаимозаменяемы: фронтенд не знает, какой из них сейчас работает,
он просто ходит на `/api/...`. Запускается всегда только один бэкенд за раз.

## Архитектура

```
                     ┌───────────────────────────┐
   :80 / :443 ──────►│           nginx           │
                     │  (HTTP или HTTPS — env)  │
                     └───────────┬───────────────┘
                     ┌───────────┴───────────────┐
                     │                           │
              location /            location /api/, /swagger-ui*, /v3/api-docs
                     │                           │
                     ▼                           ▼
             ┌───────────────┐           ┌───────────────────┐
             │   frontend    │           │      backend      │
             │ (React, nginx)│           │  crypto (Java) ИЛИ │
             │               │           │  crypto-go (Go)    │
             └───────────────┘           └─────────┬──────────┘
                                                    ▼
                                          ┌───────────────────┐
                                          │   PostgreSQL (db)  │
                                          └───────────────────┘
```

Наружу проброшен только `nginx`. `frontend`, `backend` и `db` общаются
между собой исключительно по внутренней docker-сети.

## Быстрый старт

Требуется Docker и Docker Compose v2.

```bash
cp .env.example .env
# отредактируйте .env — как минимум задайте CRYPTO_KEYSTORE_PASSWORD

# вариант на Java:
docker compose -f docker-compose.java.yml up --build -d

# вариант на Go:
docker compose -f docker-compose.go.yml up --build -d
```

Оба compose-файла читают один и тот же `.env` из текущей папки. Одновременно
поднимать оба стека не нужно (задание не требует прокидывать фронтенд в оба
бэкенда сразу — стеки переключаются по очереди). Если всё же хочется поднять
оба одновременно, задайте разные `HTTP_PORT`/`HTTPS_PORT` в отдельных `.env`
файлах и передавайте их через `--env-file`.

Погасить стек: `docker compose -f docker-compose.java.yml down` (или `.go.yml`).
Погасить с удалением данных БД и ключей: добавьте `-v`.

## Доступ после запуска

- Фронтенд: http://localhost/ (или https://localhost/, если включён HTTPS)
- Swagger UI: http://localhost/swagger-ui.html (Java) или http://localhost/swagger-ui/ (Go)
- OpenAPI JSON: http://localhost/v3/api-docs
- REST API напрямую: http://localhost/api/crypto/{encrypt,decrypt,hash,logs}

## Переменные окружения (`.env`)

См. `.env.example` — там перечислены все переменные с комментариями:

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `ENABLE_HTTPS` | `true`/`false` — включить HTTPS в nginx | `false` |
| `HTTP_PORT` | порт на хосте для HTTP | `80` |
| `HTTPS_PORT` | порт на хосте для HTTPS | `443` |
| `POSTGRES_DB` | имя базы данных | `crypto_db` |
| `POSTGRES_USER` | пользователь PostgreSQL | `postgres` |
| `POSTGRES_PASSWORD` | пароль PostgreSQL | `postgres` |
| `CRYPTO_KEYSTORE_PASSWORD` | пароль PKCS12-хранилища RSA-ключей backend'а (обязателен, без дефолта) | — |

Java- и Go-бэкенды используют одни и те же имена переменных для БД и
пароля хранилища ключей — можно переключать стеки, не меняя `.env`.

## Настройка HTTPS

По умолчанию (`ENABLE_HTTPS=false`) nginx слушает только HTTP на `HTTP_PORT`.

Чтобы включить HTTPS:

1. В `.env` поставьте `ENABLE_HTTPS=true`.
2. Перезапустите стек: `docker compose -f docker-compose.java.yml up --build -d`
   (или `.go.yml`).
3. При первом запуске с `ENABLE_HTTPS=true` контейнер `nginx` сам
   сгенерирует самоподписанный сертификат (RSA-2048, `CN=localhost`,
   срок действия 365 дней) в volume `nginx-certs` и будет переиспользовать
   его при последующих перезапусках. HTTP (`HTTP_PORT`) в этом режиме
   автоматически редиректит на HTTPS (`HTTPS_PORT`).
4. Браузер покажет предупреждение о недоверенном сертификате — это
   ожидаемо для самоподписанного сертификата, для локальной проверки
   просто подтвердите переход.

Чтобы использовать свой сертификат вместо самоподписанного — положите
`fullchain.pem` и `privkey.pem` в volume `nginx-certs` до запуска
(`docker compose ... run --rm -v ...` или `docker cp`), тогда генерация
автоматическая пропускается: скрипт проверяет наличие обоих файлов и
генерирует новый сертификат, только если их ещё нет.

## Журнал операций

Каждый вызов `/encrypt`, `/decrypt` и `/hash` сохраняется в PostgreSQL
(входные данные и результат). Посмотреть его можно на вкладке «Журнал»
во фронтенде или напрямую:

```bash
curl "http://localhost/api/crypto/logs?page=0&size=20"
curl "http://localhost/api/crypto/logs?type=HASH"
curl http://localhost/api/crypto/logs/1
```

## Локальная разработка без Docker

### Backend (Java, `crypto/`)

```bash
export CRYPTO_KEYSTORE_PASSWORD=my-strong-password
export DB_PASSWORD=postgres   # локальный PostgreSQL, БД crypto_db должна существовать
./mvnw spring-boot:run
```

Слушает `:8080`. Тесты: `./mvnw test`.

### Backend (Go, `crypto-go/`)

```bash
export CRYPTO_KEYSTORE_PASSWORD=my-strong-password
export DB_PASSWORD=postgres
go run ./cmd/server
```

Слушает `:8081` (переопределяется `SERVER_PORT`). Тесты: `go test ./...`
(интеграционные тесты репозитория/хендлеров автоматически пропускаются,
если локальный PostgreSQL недоступен).

Обоим backend'ам для локального запуска нужна локальная база `crypto_db`:

```bash
psql -U postgres -h localhost -c "CREATE DATABASE crypto_db;"
```

### Frontend (`frontend/`)

```bash
cp .env.example .env   # при необходимости поменяйте VITE_API_PROXY_TARGET
npm install
npm run dev
```

Поднимется на `http://localhost:5173`, dev-сервер Vite проксирует `/api`,
`/swagger-ui*` и `/v3/api-docs` на `VITE_API_PROXY_TARGET` (по умолчанию
`http://localhost:8080` — Java-бэкенд; для Go поставьте `:8081`).

Production-сборка (`npm run build`) собирает статику в `dist/` — именно
её раздаёт nginx внутри контейнера `frontend`; в проде фронтенд всегда
обращается по относительному `/api`, без переменных окружения.

## Структура репозитория

```
crypto/              бэкенд на Java (Spring Boot)
crypto-go/           бэкенд на Go
frontend/            веб-клиент (React + TypeScript + Vite)
nginx/               reverse-proxy с переключателем HTTP/HTTPS
docker-compose.java.yml   стек: db + crypto (Java) + frontend + nginx
docker-compose.go.yml     стек: db + crypto-go (Go) + frontend + nginx
.env.example         переменные окружения для обоих docker-compose файлов
```
