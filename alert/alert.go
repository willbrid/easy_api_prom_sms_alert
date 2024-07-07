package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"
	"easy-api-prom-alert-sms/utils"

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
		provider       config.Provider   = alertSender.config.EasyAPIPromAlertSMS.Provider
		providerParams                   = provider.Parameters
		postParams     map[string]string = map[string]string{
			providerParams.Message.ParamName: message,
		}
		queryParams utils.UrlParams
	)

	addParam := func(param config.Parameter, value string, postParams map[string]string, queryParams utils.UrlParams) error {
		switch param.ParamMethod {
		case config.PostMethod:
			postParams[param.ParamName] = value
		case config.QueryMethod:
			queryParams.AddURLParams(param.ParamName, value)
		default:
			return fmt.Errorf("bad provider parameter method: %s", param.ParamMethod)
		}

		return nil
	}

	if err := addParam(providerParams.From, providerParams.From.ParamValue, postParams, queryParams); err != nil {
		return "", "", err
	}

	if err := addParam(providerParams.To, member, postParams, queryParams); err != nil {
		return "", "", err
	}

	var encodedURL string = provider.Url
	if len(queryParams) > 0 {
		encodedURL = fmt.Sprintf("%s?%s", encodedURL, queryParams.EncodeURLParams())
	}

	body, err := utils.GetRequestBodyFromContentType(provider.ContentType, postParams)
	if err != nil {
		return "", "", err
	}

	return encodedURL, body, nil
}

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

			if err := sendSMSFromApi(url, body, simulation, provider); err != nil {
				logging.Log(logging.Error, err.Error())
			}
		}
	}

	return nil
}
