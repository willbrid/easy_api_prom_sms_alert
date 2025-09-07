package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"
	"easy-api-prom-alert-sms/utils/httpclient"

	"fmt"
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

	switch alert.Status {
	case "firing":
		message += "Started: " + alert.StartsAt.Format(time.RFC822)
	case "resolved":
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
		provider        config.Provider            = alertSender.config.EasyAPIPromAlertSMS.Provider
		providerParams                             = provider.Parameters
		httpClientParam httpclient.HttpClientParam = httpclient.HttpClientParam{
			PostParams: map[string]string{
				providerParams.Message.ParamName: strings.ReplaceAll(providerParams.Message.ParamValue, config.AlertMessageTemplate, message),
			},
			QueryParams: map[string]string{},
		}
	)

	httpClientParam.AddParam(providerParams.From.ParamMethod, providerParams.From.ParamName, providerParams.From.ParamValue)
	httpClientParam.AddParam(providerParams.To.ParamMethod, providerParams.To.ParamName, member)
	if len(providerParams.ExtraParams) > 0 {
		for _, extraParam := range providerParams.ExtraParams {
			httpClientParam.AddParam(extraParam.ParamMethod, extraParam.ParamName, extraParam.ParamValue)
		}
	}

	var encodedURL string = provider.Url
	if len(httpClientParam.QueryParams) > 0 {
		encodedURL = fmt.Sprintf("%s?%s", encodedURL, httpClientParam.EncodeQueryParams())
	}

	body, err := httpClientParam.EncodePostParams(provider.ContentType)
	if err != nil {
		return "", "", err
	}

	return encodedURL, body, nil
}

// sendAlert help to send alert with alertsender object
func (alertSender *AlertSender) sendAlert() error {
	for _, alert := range alertSender.data.Alerts {
		alertMsg := alertSender.getMsgFromAlert(alert)
		recipientName := alertSender.getRecipientFromAlert(alert)
		members := alertSender.getRecipientMembers(recipientName)
		provider := alertSender.config.EasyAPIPromAlertSMS.Provider
		simulation := alertSender.config.EasyAPIPromAlertSMS.Simulation

		for _, member := range members {
			url, body, err := alertSender.getUrlAndBody(member, alertMsg)

			if err != nil {
				return err
			}

			go func() {
				if err := sendSMSFromApi(url, body, simulation, provider); err != nil {
					logging.Log(logging.Error, "%s", err.Error())
				}
			}()
		}
	}

	return nil
}
