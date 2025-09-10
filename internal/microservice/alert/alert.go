package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/pkg/httpclient"
	"easy-api-prom-alert-sms/pkg/httpclientparam"

	"fmt"
	"strings"
)

type AlertMicroservice struct {
	Provider config.Provider
}

func NewAlertMicroservice(provider config.Provider) *AlertMicroservice {
	return &AlertMicroservice{Provider: provider}
}

func (ams *AlertMicroservice) Consume(url string, body string) error {
	provider := ams.Provider
	var headers map[string]string = map[string]string{
		"Content-Type": provider.ContentType,
	}

	if provider.Authentication.Enabled {
		headers["Authorization"] = fmt.Sprintf("%s %s", provider.Authentication.AuthorizationType, provider.Authentication.AuthorizationCredential)
	}

	_, err := httpclient.Post(url, strings.NewReader(body), httpclient.Options{
		Headers:            headers,
		Timeout:            provider.Timeout,
		InsecureSkipVerify: provider.InsecureSkipVerify,
	})

	if err != nil {
		return err
	}

	return nil
}

func (ams *AlertMicroservice) GetUrlAndBody(dest string, message string) (string, string, error) {
	httpClientParam := httpclientparam.NewHttpClientParam()
	providerParams := ams.Provider.Parameters
	httpClientParam.PostParams[providerParams.Message.ParamName] = strings.ReplaceAll(providerParams.Message.ParamValue, config.AlertMessageTemplate, message)

	httpClientParam.AddParam(providerParams.From.ParamMethod, providerParams.From.ParamName, providerParams.From.ParamValue)
	httpClientParam.AddParam(providerParams.To.ParamMethod, providerParams.To.ParamName, dest)

	if len(providerParams.ExtraParams) > 0 {
		for _, extraParam := range providerParams.ExtraParams {
			httpClientParam.AddParam(extraParam.ParamMethod, extraParam.ParamName, extraParam.ParamValue)
		}
	}
	var encodedURL string = ams.Provider.Url
	if len(httpClientParam.QueryParams) > 0 {
		encodedURL = fmt.Sprintf("%s?%s", encodedURL, httpClientParam.EncodeQueryParams())
	}

	body, err := httpClientParam.EncodePostParams(ams.Provider.ContentType)
	if err != nil {
		return "", "", err
	}

	return encodedURL, body, nil
}
