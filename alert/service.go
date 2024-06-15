package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"
	"easy-api-prom-alert-sms/utils"

	"fmt"
	"strings"
)

func sendSMSFromApi(url string, body string, simulation bool, provider config.Provider) error {
	if simulation {
		logging.Log(logging.Info, "send request with url %s and body %s", url, string(body))
		return nil
	}

	var headers map[string]string = map[string]string{
		"Content-Type": provider.ContentType,
	}
	if provider.Authentication.Enabled {
		headers["Authorization"] = fmt.Sprintf("%s %s", provider.Authentication.AuthorizationType, provider.Authentication.AuthorizationCredential)
	}

	if err := utils.Post(url, strings.NewReader(body), utils.Options{
		Headers: headers,
		Timeout: provider.Timeout,
	}); err != nil {
		return err
	}

	return nil
}
