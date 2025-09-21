package domain

import (
	"easy-api-prom-alert-sms/config"

	"github.com/prometheus/alertmanager/template"
)

type Alert struct {
	Data *template.Data
}

type AlertConfig struct {
	Recipients           []config.Recipient
	DefaultRecipientName string
	Simulation           bool
}
