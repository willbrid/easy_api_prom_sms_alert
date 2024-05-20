package alert

import (
	"easy-api-prom-alert-sms/config"

	"github.com/prometheus/alertmanager/template"
)

type AlertSender struct {
	config *config.Config
	data   *template.Data
}

func NewAlertSender(config *config.Config) *AlertSender {
	return &AlertSender{config, nil}
}

func (alertSender *AlertSender) setData(data *template.Data) {
	alertSender.data = data
}
