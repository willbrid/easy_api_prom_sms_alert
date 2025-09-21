#!/bin/sh

exec /usr/local/bin/easy-api-prom-sms-alert \
    --config-file $EASY_API_PROM_SMS_ALERT_CONFIG_FILE \
    --port $EASY_API_PROM_SMS_ALERT_PORT \
    --enable-https $EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS \
    --cert-file $EASY_API_PROM_SMS_ALERT_CERT_FILE \
    --key-file $EASY_API_PROM_SMS_ALERT_KEY_FILE