# syntax=docker/dockerfile:1
FROM nginx:1.27-alpine

RUN apk add --no-cache openssl

COPY conf/http.conf conf/https.conf conf/proxy-common.conf /etc/nginx/conf.available/
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && rm -f /etc/nginx/conf.d/default.conf

ENV ENABLE_HTTPS=false

EXPOSE 80 443
ENTRYPOINT ["/entrypoint.sh"]
