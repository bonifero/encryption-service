# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/crypto-go ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /out/crypto-go .

RUN mkdir -p /app/keystore

EXPOSE 8081
ENTRYPOINT ["./crypto-go"]
