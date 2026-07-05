#!/bin/sh
set -e

CERT_DIR=/etc/nginx/certs
mkdir -p "$CERT_DIR"

if [ "$ENABLE_HTTPS" = "true" ]; then
  if [ ! -f "$CERT_DIR/fullchain.pem" ] || [ ! -f "$CERT_DIR/privkey.pem" ]; then
    echo "ENABLE_HTTPS=true: сертификат не найден в $CERT_DIR, генерирую самоподписанный (только для разработки/демо)"
    openssl req -x509 -nodes -newkey rsa:2048 \
      -keyout "$CERT_DIR/privkey.pem" \
      -out "$CERT_DIR/fullchain.pem" \
      -days 365 \
      -subj "/CN=localhost/O=crypto-service/C=RU" \
      -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
  else
    echo "ENABLE_HTTPS=true: использую существующий сертификат из $CERT_DIR"
  fi
  cp /etc/nginx/conf.available/https.conf /etc/nginx/conf.d/default.conf
else
  echo "ENABLE_HTTPS=false: работаю по обычному HTTP"
  cp /etc/nginx/conf.available/http.conf /etc/nginx/conf.d/default.conf
fi

exec nginx -g "daemon off;"
