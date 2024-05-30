package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

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

func (alertSender *AlertSender) getBodyAndUrl(member string, message string) (map[string]string, string) {
	postParams := map[string]string{
		alertSender.config.EasyAPIPromAlertSMS.Provider.Parameters.Message.ParamName: message,
	}
	queryParams := url.Values{}

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

	encodedURL := fmt.Sprintf("%s?%s", alertSender.config.EasyAPIPromAlertSMS.Provider.Url, queryParams.Encode())
	return postParams, encodedURL
}

func (alertSender *AlertSender) sendAlert() error {
	for _, alert := range alertSender.data.Alerts {
		alertMsg := alertSender.getMsgFromAlert(alert)
		recipientName := alertSender.getRecipientFromAlert(alert)
		members := alertSender.getRecipientMembers(recipientName)

		for _, member := range members {
			var builder strings.Builder
			body, encodedURL := alertSender.getBodyAndUrl(member, alertMsg)
			if err := json.NewEncoder(&builder).Encode(body); err != nil {
				return err
			}

			if err := sendSMSFromProviderApi(
				encodedURL,
				builder.String(),
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

func sendSMSFromProviderApi(encodedURL string, body string, authEnable bool, authType string, authCred string, timeout time.Duration, simulation bool) error {
	if simulation {
		logging.Log(logging.Info, "successful send request with url %s and body %s", encodedURL, body)
	} else {
		client := &http.Client{
			Timeout: timeout,
		}

		req, err := http.NewRequest("POST", encodedURL, strings.NewReader(body))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		if authEnable {
			req.Header.Set("Authorization", authType+" "+authCred)
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		var respBody []byte
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			logging.Log(logging.Error, "Failed to read response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("request failed with status : %s", resp.Status)
		}

		logging.Log(logging.Info, "successful send request with url %s and body %s", encodedURL, body)
		logging.Log(logging.Info, "response body %s", string(respBody))
	}

	return nil
}
