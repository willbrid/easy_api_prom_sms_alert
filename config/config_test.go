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

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := "the field auth.username is required and must be a string between 2 and 25 characters long"

			if err == nil {
				t.Fatalf("no error returned, expected:\n%v", expected)
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
    username: "test"
    password: ""
`,
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxx
`,
		`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: testxxx
`,
	}

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := "the field auth.password is required and must be a string between 8 and 255 characters long"

			if err == nil {
				t.Fatalf("no error returned, expected:\n%v", expected)
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

	for index, configContent := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig #%v", index), func(subT *testing.T) {
			filename, err := utils.CreateConfigFileForTesting(configContent)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer os.Remove(filename)

			_, err = LoadConfig(filename, validate)

			expected := "the field provider.url is required and must be a valid url"

			if err == nil {
				t.Fatalf("no error returned, expected:\n%v", expected)
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
      authorization:
        type: ''
        credential: ''
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
      authorization:
        type: 'xxxxx'
        credential: ''
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
      authorization:
        type: 'Bearer'
        credential: ''
`,
	}

	expectations := []string{
		"when provider.authentication is enabled, you should provider authorization config",
		"when provider.authentication.authorization is used, the field provider.authentication.authorization.type must be among the values : Bearer, Basic, ApiKey",
		"when provider.authentication.authorization is used, the field provider.authentication.authorization.credential is required and must be a string at most 255 characters long",
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
				t.Fatalf("no error returned, expected:\n%v", expected)
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
    parameters:
      from:
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
		"the field provider.parameters.from.param_name is required and must be a string at most 25 characters long",
		"the field provider.parameters.from.param_method must be among the values : post and query",
		"the field provider.parameters.from.param_value is required and must be a string at most 25 characters long",
		"the field provider.parameters.to.param_name is required and must be a string at most 25 characters long",
		"the field provider.parameters.to.param_method must be among the values : post and query",
		"the field provider.parameters.to.param_value is required and must be a string at most 25 characters long",
		"the field provider.parameters.message.param_name is required and must be a string at most 25 characters long",
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
				t.Fatalf("no error returned, expected:\n%v", expected)
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
      authorization:
        type: "Bearer"
        credential: "xxxxx"
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
		"error converting recipients to slice of interface{}",
		"the field recipients[].name is required and must be a string at most 25 characters long",
		"error converting recipientMap['members'] to slice of interface{}",
		"the field recipients[].members[] is required and must be a string at most 25 characters long",
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
				t.Fatalf("no error returned, expected:\n%v", expected)
			}

			if err.Error() != expected {
				t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
			}
		})
	}
}
