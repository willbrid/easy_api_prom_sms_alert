package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/alertmanager/template"
)

func (alertSender *AlertSender) AlertHandler(resp http.ResponseWriter, req *http.Request) {
	var alertData template.Data

	if err := json.NewDecoder(req.Body).Decode(&alertData); err != nil {
		logging.Log(logging.Error, "failed to parse content : %s", err.Error())
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	alertSender.setData(&alertData)

	go func() {
		if err := alertSender.sendAlert(); err != nil {
			logging.Log(logging.Error, "failed to send alert : %s", err.Error())
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
	}()

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNoContent)
}

func (alertSender *AlertSender) sendAlert() error {
	for _, alert := range alertSender.data.Alerts {
		alertMsg := alertSender.getMsgFromAlert(alert)
		recipientName := alertSender.getRecipientFromAlert(alert)
		members := alertSender.getRecipientMembers(recipientName)

		for _, member := range members {
			var builder strings.Builder
			body := map[string]string{
				alertSender.config.EasyAPIPromAlertSMS.Provider.From:    alertSender.config.EasyAPIPromAlertSMS.Provider.FromValue,
				alertSender.config.EasyAPIPromAlertSMS.Provider.To:      member,
				alertSender.config.EasyAPIPromAlertSMS.Provider.Message: alertMsg,
			}
			if err := json.NewEncoder(&builder).Encode(body); err != nil {
				return err
			}

			if alertSender.config.EasyAPIPromAlertSMS.Simulation {
				logging.Log(logging.Info, builder.String())
			} else {
				if err := consumeProviderApi(alertSender.config, builder.String()); err != nil {
					logging.Log(logging.Error, err.Error())
					continue
				}
			}
		}
	}

	return nil
}

func consumeProviderApi(config *config.Config, message string) error {
	client := &http.Client{
		Timeout: config.EasyAPIPromAlertSMS.Provider.Timeout,
	}

	req, err := http.NewRequest("POST", config.EasyAPIPromAlertSMS.Provider.Url, strings.NewReader(message))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if config.EasyAPIPromAlertSMS.Provider.Authentication.Enabled {
		req.Header.Set("Authorization", config.EasyAPIPromAlertSMS.Provider.Authentication.Authorization.Type+" "+config.EasyAPIPromAlertSMS.Provider.Authentication.Authorization.Credential)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status : %s", resp.Status)
	}

	return nil
}
