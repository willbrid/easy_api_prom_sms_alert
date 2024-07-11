package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"
	"easy-api-prom-alert-sms/utils/httpclient"

	"fmt"
	"strings"
)

func sendSMSFromApi(url string, body string, simulation bool, provider config.Provider) error {
	if simulation {
		logging.Log(logging.Info, "send request with url %s and body %s", url, body)
		return nil
	}

	var headers map[string]string = map[string]string{
		"Content-Type": provider.ContentType,
	}
	if provider.Authentication.Enabled {
		headers["Authorization"] = fmt.Sprintf("%s %s", provider.Authentication.AuthorizationType, provider.Authentication.AuthorizationCredential)
	}

	if err := httpclient.Post(url, strings.NewReader(body), httpclient.Options{
		Headers: headers,
		Timeout: provider.Timeout,
	}); err != nil {
		return err
	}

	return nil
}
