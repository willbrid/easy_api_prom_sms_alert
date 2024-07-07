package utils

import "fmt"

type UrlParams map[string]string

// EncodeURLParams help to encode url params
func (urlParams UrlParams) EncodeURLParams() string {
	var encodedURLParams string = ""

	for key, value := range urlParams {
		encodedURLParams += fmt.Sprintf("&%s=%s", key, value)
	}
	encodedURLParams = encodedURLParams[1:]

	return encodedURLParams
}

// AddURLParam help to add a new url param
func (urlParams UrlParams) AddURLParams(key string, value string) {
	urlParams[key] = value
}
