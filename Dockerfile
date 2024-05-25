FROM golang:1.21 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o easy-api-prom-sms-alert .

FROM alpine:3.20

COPY --from=builder /build/easy-api-prom-sms-alert /usr/local/bin

RUN mkdir /etc/easy-api-prom-sms-alert
COPY --from=builder --chown=nobody /build/config.default.yaml /etc/easy-api-prom-sms-alert/config.yaml

RUN apk update && apk add --no-cache ca-certificates

ENV EASY_API_PROM_SMS_ALERT_PORT=5957 

USER nobody
EXPOSE $EASY_API_PROM_SMS_ALERT_PORT

ENTRYPOINT ["easy-api-prom-sms-alert"]
CMD ["--config-file", "/etc/easy-api-prom-sms-alert/config.yaml", "--port", $EASY_API_PROM_SMS_ALERT_PORT]