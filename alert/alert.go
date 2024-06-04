package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"
	"easy-api-prom-alert-sms/utils"

	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

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
		recipientName = alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.To.ParamValue
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

	if alert.Status == "firing" {
		message += "Started: " + alert.StartsAt.Format(time.RFC822)
	} else if alert.Status == "resolved" {
		message += "Resolved: " + alert.EndsAt.Format(time.RFC822)
	}

	return message
}

// getRecipientMembers get recipient members from recipient name
func (alertSender *AlertSender) getRecipientMembers(recipientName string) []string {
	recipients := alertSender.config.EasyAPIPromAlertSMS.Recipients
	var recipient config.Recipient

	for _, recipientItem := range recipients {
		if recipientItem.Name == recipientName {
			recipient = recipientItem
			break
		}
	}

	return recipient.Members
}

// getUrlAndBody help to get parsed url and body
func (alertSender *AlertSender) getUrlAndBody(member string, message string) (string, string, error) {
	var (
		postParams map[string]string = map[string]string{
			alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.Message.ParamName: message,
		}
		queryParams url.Values = url.Values{}
	)

	if alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.From.ParamMethod == config.PostMethod {
		postParams[alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.From.ParamName] = alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.From.ParamValue
	} else {
		queryParams.Add(alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.From.ParamName, alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.From.ParamValue)
	}

	if alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.To.ParamMethod == config.PostMethod {
		postParams[alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.To.ParamName] = member
	} else {
		queryParams.Add(alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.To.ParamName, member)
	}

	var (
		builder    strings.Builder
		encodedURL string
	)

	if len(queryParams) > 0 {
		encodedURL = fmt.Sprintf("%s?%s", alertSender.config.EasyAPIPromAlertSMS.Provider.Url, queryParams.Encode())
	} else {
		encodedURL = alertSender.config.EasyAPIPromAlertSMS.Provider.Url
	}

	if err := json.NewEncoder(&builder).Encode(postParams); err != nil {
		return "", "", err
	}

	return encodedURL, builder.String(), nil
}

func (alertSender *AlertSender) sendAlert() error {
	for _, alert := range alertSender.data.Alerts {
		alertMsg := alertSender.getMsgFromAlert(alert)
		recipientName := alertSender.getRecipientFromAlert(alert)
		members := alertSender.getRecipientMembers(recipientName)

		for _, member := range members {
			url, body, err := alertSender.getUrlAndBody(member, alertMsg)

			if err != nil {
				return err
			}

			if err := utils.SendSMSFromApi(
				url,
				body,
				alertSender.config.EasyAPIPromAlertSMS.Provider.Authentication.Enabled,
				alertSender.config.EasyAPIPromAlertSMS.Provider.Authentication.Type,
				alertSender.config.EasyAPIPromAlertSMS.Provider.Authentication.Credential,
				alertSender.config.EasyAPIPromAlertSMS.Provider.Timeout,
				alertSender.config.EasyAPIPromAlertSMS.Simulation,
			); err != nil {
				logging.Log(logging.Error, err.Error())
			}
		}
	}

	return nil
}
