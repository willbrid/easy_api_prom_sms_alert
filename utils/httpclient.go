package utils

import (
	"easy-api-prom-alert-sms/logging"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Options struct {
	Headers map[string]string
	Timeout time.Duration
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

func GetRequestBodyFromContentType(contentType string, postParams map[string]string) (string, error) {
	var reqBody string

	switch contentType {
	case "application/x-www-form-urlencoded":
		reqBody = ""
		for key, value := range postParams {
			reqBody += fmt.Sprintf("&%s=%s", key, value)
		}
		reqBody = reqBody[1:]

	case "application/json":
		postParamStr, err := json.Marshal(postParams)
		if err != nil {
			return "", err
		}
		reqBody = string(postParamStr)

	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	return reqBody, nil
}
