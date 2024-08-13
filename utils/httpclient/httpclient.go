package httpclient

import (
	"crypto/tls"
	"easy-api-prom-alert-sms/logging"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	PostMethod  = "post"
	QueryMethod = "query"
)

type HttpClientParam struct {
	PostParams  map[string]string
	QueryParams map[string]string
}

type Options struct {
	Headers            map[string]string
	Timeout            time.Duration
	InsecureSkipVerify bool
}

var httpClient *http.Client

const (
	MaxIdleConnections int = 20
)

func init() {
	httpClient = createHTTPClient()
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
	}

	return client
}

func Post(url string, body io.Reader, options Options) error {
	httpClient.Timeout = options.Timeout
	transport := httpClient.Transport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: options.InsecureSkipVerify}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	for headerKey, headerValue := range options.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var respBody []byte
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		logging.Log(logging.Error, "failed to read response body : %v", err)
	}

	bodyStr, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	logging.Log(logging.Info, "send request with url %s and body %s", url, string(bodyStr))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request failed with status : %s and response body : %s", resp.Status, string(respBody))
	}

	logging.Log(logging.Info, "response body %s", string(respBody))
	return nil
}

func (httpClientParam HttpClientParam) AddParam(method, paramKey, paramValue string) {
	switch method {
	case PostMethod:
		httpClientParam.PostParams[paramKey] = paramValue
	case QueryMethod:
		httpClientParam.QueryParams[paramKey] = paramValue
	default:
		httpClientParam.PostParams[paramKey] = paramValue
	}
}

func (httpClientParam HttpClientParam) EncodeQueryParams() string {
	return parseUrlParams(httpClientParam.QueryParams)
}

func (httpClientParam HttpClientParam) EncodePostParams(contentType string) (string, error) {
	var postParams map[string]string = httpClientParam.PostParams
	var encodedPostParams string

	switch contentType {
	case "application/x-www-form-urlencoded":
		encodedPostParams = parseUrlParams(postParams)

	case "application/json":
		postParamStr, err := json.Marshal(postParams)
		if err != nil {
			return "", err
		}
		encodedPostParams = string(postParamStr)
	}

	return encodedPostParams, nil
}

func parseUrlParams(params map[string]string) string {
	var encodedUrlParams string

	for key, value := range params {
		encodedUrlParams += fmt.Sprintf("&%s=%s", key, value)
	}
	encodedUrlParams = encodedUrlParams[1:]

	return encodedUrlParams
}
