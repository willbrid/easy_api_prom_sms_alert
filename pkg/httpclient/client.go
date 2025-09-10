package httpclient

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

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

func setHTTPClientOptions(options Options) {
	httpClient.Timeout = options.Timeout
	transport := httpClient.Transport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: options.InsecureSkipVerify}
}

func Post(url string, body io.Reader, options Options) (*http.Response, error) {
	setHTTPClientOptions(options)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	return httpClient.Do(req)
}
