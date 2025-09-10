package util

import (
	"sort"
	"strings"
	"time"

	"github.com/prometheus/alertmanager/template"
)

// GetRecipientFromAlert get recipient name from template.Alert
func GetRecipientFromAlert(alert template.Alert, defaultRecipient string) string {
	var recipientName string = defaultRecipient

	if value, exists := alert.Labels["team"]; exists {
		recipientName = value
	}

	return recipientName
}

// GetMsgFromAlert generate message to send from alert
func GetMsgFromAlert(alert template.Alert) string {
	var (
		pairs   []string = []string{}
		message string
	)

	for k, v := range alert.Labels {
		if k != "team" {
			pairs = append(pairs, k+": "+v)
		}
	}
	sort.Strings(pairs)
	message = strings.ToUpper(alert.Status) + "\n" + strings.Join(pairs, "\n") + "\n"

	if summary, exists := alert.Annotations["summary"]; exists && summary != "" {
		message += "summary: " + summary + "\n"
	}

	if description, exists := alert.Annotations["description"]; exists && description != "" {
		message += "description: " + description + "\n"
	}

	switch alert.Status {
	case "firing":
		message += "Started: " + alert.StartsAt.Format(time.RFC822)
	case "resolved":
		message += "Resolved: " + alert.EndsAt.Format(time.RFC822)
	}

	return message
}
