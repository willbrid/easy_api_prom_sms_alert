package config

import (
	"easy-api-prom-alert-sms/utils"

	"fmt"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func TestLoadConfigFileNotFound(t *testing.T) {
	var filename string

	_, err := LoadConfig(filename, validate)

	expected := "Config File \"config\" Not Found in \"[]\""

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestLoadConfigFileNotExist(t *testing.T) {
	var filename string = "nonexistentfile.yaml"

	_, err := LoadConfig(filename, validate)

	expected := "open nonexistentfile.yaml: no such file or directory"

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestUsernameFieldWithAuthEnabled(t *testing.T) {
	configSlices := []string{
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: ""
`,
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "x"
`,
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`,
	}

	expectations := []string{
		"validation failed on field 'Username' for condition 'required_if'",
		"validation failed on field 'Username' for condition 'min'",
		"validation failed on field 'Username' for condition 'max'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}

func TestPasswordFieldWithAuthEnabled(t *testing.T) {
	configSlices := []string{
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: ""
`,
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxx
`,
	}

	expectations := []string{
		"validation failed on field 'Password' for condition 'required_if'",
		"validation failed on field 'Password' for condition 'min'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}

func TestProviderUrl(t *testing.T) {
	configSlices := []string{
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: ""
`,
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://"
`,
	}

	expectations := []string{
		"validation failed on field 'Url' for condition 'required'",
		"validation failed on field 'Url' for condition 'url'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}

func TestProviderAuthField(t *testing.T) {
	configSlices := []string{
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
`,
		`---
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
`,
		`---
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
`,
		`---
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
`,
	}

	expectations := []string{
		"validation failed on field 'AuthorizationType' for condition 'required_if'",
		"validation failed on field 'AuthorizationType' for condition 'required_if'",
		"validation failed on field 'AuthorizationCredential' for condition 'required_if'",
		"validation failed on field 'AuthorizationType' for condition 'required_if'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}

func TestProviderParamField(t *testing.T) {
	configSlices := []string{
		`---
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
`,
		`---
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
`,
		`---
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
`,
		`---
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
`,
		`---
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
`,
		`---
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
`,
		`---
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
`,
		`---
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
        param_value: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`,
		`---
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
`,
	}

	expectations := []string{
		"validation failed on field 'ParamValue' for condition 'required'",
		"validation failed on field 'From' for condition 'required'",
		"validation failed on field 'ParamName' for condition 'max'",
		"validation failed on field 'ParamValue' for condition 'required'",
		"validation failed on field 'ParamMethod' for condition 'oneof'",
		"validation failed on field 'ParamName' for condition 'max'",
		"validation failed on field 'ParamValue' for condition 'required'",
		"validation failed on field 'ParamValue' for condition 'max'",
		"validation failed on field 'ParamName' for condition 'max'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}

func TestRecipientField(t *testing.T) {
	configSlices := []string{
		`---
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
`,
		`---
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
`,
		`---
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
`,
		`---
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
`,
	}

	expectations := []string{
		"validation failed on field 'Recipients' for condition 'gt'",
		"validation failed on field 'Name' for condition 'required'",
		"validation failed on field 'Members[0]' for condition 'min'",
		"validation failed on field 'Members[0]' for condition 'max'",
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := expectations[index]

			if err == nil {
				t.Errorf("no error returned, expected:\n%v", expected)
			}

			if err != nil && err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}
