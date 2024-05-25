package alert

import (
	"easy-api-prom-alert-sms/config"

	"sort"
	"strings"

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

// getRecipientFromAlert get recipient name from alert
func (alertSender *AlertSender) getRecipientFromAlert(alert template.Alert) string {
	var recipientName string

	if value, exists := alert.Labels["team"]; exists {
		recipientName = value
	} else {
		recipientName = alertSender.config.EasyAPIPromAlertSMS.Recipients[0].Name
	}

	return recipientName
}

// getMsgFromAlert generate message to send from alert
func (alertSender *AlertSender) getMsgFromAlert(alert template.Alert) string {
	var (
		pairs   []string = []string{}
		message string
	)

	for k, v := range alert.Labels {
		if k != "team" {
			pairs = append(pairs, k+"= "+v)
		}
	}
	sort.Strings(pairs)
	message = strings.ToUpper(alert.Status) + "\n" + strings.Join(pairs, "\n")

	if summary, exists := alert.Annotations["summary"]; exists && summary != "" {
		message += "summary: " + summary + "\n"
	}

	if description, exists := alert.Annotations["description"]; exists && description != "" {
		message += "description: " + description + "\n"
	}

	return message
}

// getRecipientMembers get recipient members from recipient name
func (alertSender *AlertSender) getRecipientMembers(recipientName string) []string {
	recipients := alertSender.config.EasyAPIPromAlertSMS.Recipients
	var recipient config.Recipient

	for _, recipientItem := range recipients {
		if recipient.Name == recipientName {
			recipient = recipientItem
			break
		}
	}

	return recipient.Members
}
