package config_test

import (
	"easy-api-prom-alert-sms/config"

	"bytes"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func triggerTest(t *testing.T, yamlConfig []byte, expectations []string, index int) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer([]byte(yamlConfig))); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	_, err := config.LoadConfig(v, validate)

	expected := expectations[index]

	if err == nil {
		t.Errorf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestReadConfigFile_ReturnFileNotFoundError(t *testing.T) {
	t.Parallel()

	var filename string

	_, err := config.ReadConfigFile(filename)
	expected := "configuration file '' not found"

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestReadConfigFile_ReturnFileNotExistError(t *testing.T) {
	t.Parallel()

	var filename string = "nonexistentfile.yaml"

	_, err := config.ReadConfigFile(filename)

	expected := "open nonexistentfile.yaml: no such file or directory"

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestLoadConfig_ReturnErrorWithBadUsernameInputWhenAuthEnabled(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "x"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
	}

	expectations := []string{
		"validation failed on field 'Username' for condition 'required_if'",
		"validation failed on field 'Username' for condition 'min'",
		"validation failed on field 'Username' for condition 'max'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadPasswordInputWhenAuthEnabled(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxx
`),
	}

	expectations := []string{
		"validation failed on field 'Password' for condition 'required_if'",
		"validation failed on field 'Password' for condition 'min'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadProviderUrl(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    content_type: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    content_type: "xxxxx"
`),
	}

	expectations := []string{
		"validation failed on field 'Url' for condition 'required'",
		"validation failed on field 'Url' for condition 'url'",
		"validation failed on field 'ContentType' for condition 'required'",
		"validation failed on field 'ContentType' for condition 'oneof'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadProviderAuthInput(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: ''
      authorization_credential: ''
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: 'xxxxx'
      authorization_credential: ''
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: ''
      authorization_credential: 'xxxxxxxxxx'
`),
	}

	expectations := []string{
		"validation failed on field 'AuthorizationType' for condition 'required_if'",
		"validation failed on field 'AuthorizationType' for condition 'required_if'",
		"validation failed on field 'AuthorizationCredential' for condition 'required_if'",
		"validation failed on field 'AuthorizationType' for condition 'required_if'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadProviderParamInput(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_value: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_value: "xxxxx"
        param_method: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
	}

	expectations := []string{
		"validation failed on field 'ParamValue' for condition 'required'",
		"validation failed on field 'ParamValue' for condition 'required'",
		"validation failed on field 'ParamName' for condition 'max'",
		"validation failed on field 'ParamValue' for condition 'required'",
		"validation failed on field 'ParamMethod' for condition 'oneof'",
		"validation failed on field 'ParamName' for condition 'max'",
		"validation failed on field 'ParamValue' for condition 'required'",
		"validation failed on field 'ParamName' for condition 'max'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadRecipientInput(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
  recipients:
  - name: ""
    members:
    - "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
  recipients:
  - name: "admin"
    members:
    - ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
  recipients:
  - name: "admin"
    members:
    - "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
	}

	expectations := []string{
		"validation failed on field 'Recipients' for condition 'gt'",
		"validation failed on field 'Name' for condition 'required'",
		"validation failed on field 'Members[0]' for condition 'min'",
		"validation failed on field 'Members[0]' for condition 'max'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}
