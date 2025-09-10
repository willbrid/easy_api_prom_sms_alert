FROM golang:1.24 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o easy-api-prom-sms-alert ./cmd

FROM alpine:3.20

COPY --from=builder /build/easy-api-prom-sms-alert /usr/local/bin
COPY --from=builder /build/entrypoint.sh /usr/local/bin
RUN chmod +x /usr/local/bin/entrypoint.sh

RUN mkdir /etc/easy-api-prom-sms-alert
RUN mkdir /etc/easy-api-prom-sms-alert/tls
COPY --from=builder --chown=nobody /build/fixtures/config.default.yaml /etc/easy-api-prom-sms-alert/config.yaml
COPY --from=builder --chown=nobody /build/fixtures/tls/server.crt /etc/easy-api-prom-sms-alert/tls/server.crt
COPY --from=builder --chown=nobody /build/fixtures/tls/server.key /etc/easy-api-prom-sms-alert/tls/server.key

RUN apk update && apk add --no-cache ca-certificates

ENV EASY_API_PROM_SMS_ALERT_PORT=5957
ENV EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS="true"

USER nobody
EXPOSE $EASY_API_PROM_SMS_ALERT_PORT

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]