#!/bin/sh

exec /usr/local/bin/easy-api-prom-sms-alert \
    --config-file /etc/easy-api-prom-sms-alert/config.yaml \
    --port $EASY_API_PROM_SMS_ALERT_PORT \
    --enable-https $EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS \
    --cert-file /etc/easy-api-prom-sms-alert/tls/server.crt \
    --key-file /etc/easy-api-prom-sms-alert/tls/server.key