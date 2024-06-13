package utils

import (
	"easy-api-prom-alert-sms/logging"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetRequestBodyFromContentType(contentType string, postParams map[string]string) (io.Reader, error) {
	var reqBody io.Reader

	switch contentType {
	case "application/x-www-form-urlencoded":
		data := url.Values{}
		for key, value := range postParams {
			data.Set(key, value)
		}
		reqBody = strings.NewReader(data.Encode())

	case "application/json":
		postParamStr, err := json.Marshal(postParams)
		if err != nil {
			return nil, err
		}
		reqBody = strings.NewReader(string(postParamStr))

	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return reqBody, nil
}

// SendSMSFromApi send sms through an api specification
func SendSMSFromApi(url string, body io.Reader, contentType string, authEnable bool, authType string, authCred string, timeout time.Duration, simulation bool) error {
	bodyStr, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	if simulation {
		logging.Log(logging.Info, "successful send request with url %s and body %s", url, string(bodyStr))
		return nil
	}

	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", contentType)
	if authEnable {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", authType, authCred))
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

	logging.Log(logging.Info, "successful send request with url %s and body %s", url, string(bodyStr))
	logging.Log(logging.Info, "response body %s", string(respBody))

	return nil
}
