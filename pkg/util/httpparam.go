package util

import (
	"encoding/json"
	"fmt"
)

const (
	PostMethod  = "post"
	QueryMethod = "query"
)

type HttpClientParam struct {
	PostParams  map[string]string
	QueryParams map[string]string
}

func NewHttpClientParam() *HttpClientParam {
	return &HttpClientParam{
		PostParams:  map[string]string{},
		QueryParams: map[string]string{},
	}
}

// AddParam adds a parameter to the appropriate map based on the method type.
func (httpClientParam *HttpClientParam) AddParam(method, paramKey, paramValue string) {
	switch method {
	case PostMethod:
		httpClientParam.PostParams[paramKey] = paramValue
	case QueryMethod:
		httpClientParam.QueryParams[paramKey] = paramValue
	default:
		httpClientParam.PostParams[paramKey] = paramValue
	}
}

// EncodeQueryParams encodes the query parameters into a URL-encoded string.
func (httpClientParam *HttpClientParam) EncodeQueryParams() string {
	return parseUrlParams(httpClientParam.QueryParams)
}

func (httpClientParam *HttpClientParam) EncodePostParams(contentType string) (string, error) {
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

// parseUrlParams encodes URL parameters from a map to a query string format.
func parseUrlParams(params map[string]string) string {
	var encodedUrlParams string

	for key, value := range params {
		encodedUrlParams += fmt.Sprintf("&%s=%s", key, value)
	}
	encodedUrlParams = encodedUrlParams[1:]

	return encodedUrlParams
}
