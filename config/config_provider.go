package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Provider struct {
	Url             string        `mapstructure:"url"`
	Timeout         time.Duration `mapstructure:"timeout"`
	*Authentication `mapstructure:"authentication"`
	Parameters      struct {
		From    Parameter `mapstructure:"from"`
		To      Parameter `mapstructure:"to"`
		Message Parameter `mapstructure:"message"`
	} `mapstructure:"parameters"`
}

type Authentication struct {
	Enabled        bool `mapstructure:"enabled"`
	*Authorization `mapstructure:"authorization"`
}

type Authorization struct {
	Type       string `mapstructure:"type"`
	Credential string `mapstructure:"credential"`
}

const (
	PostMethod   string = "post"
	QueryMethod  string = "query"
	NoParamName  string = "none"
	NoParamValue string = "none"
)

type Parameter struct {
	ParamName   string `mapstructure:"param_name"`
	ParamValue  string `mapstructure:"param_value"`
	ParamMethod string `mapstructure:"param_method"`
}

// validateProviderParameter validate a specific provider parameter : from, to and message
func validateProviderParameter(v *viper.Viper, validate *validator.Validate, paramKey string) error {
	providerParamField := "easy_api_prom_sms_alert.provider.parameters." + paramKey
	providerParamName := v.GetString(providerParamField + ".param_name")
	providerParamValue := v.GetString(providerParamField + ".param_value")
	providerParamMethod := v.GetString(providerParamField + ".param_method")

	if err := validate.Var(providerParamName, "required,max=25"); err != nil {
		return fmt.Errorf("the field provider.parameters.%s.param_name is required and must be a string at most 25 characters long", paramKey)
	}

	if err := validate.Var(providerParamMethod, "required,oneof="+PostMethod+" "+QueryMethod); err != nil {
		return fmt.Errorf("the field provider.parameters.%s.param_method must be among the values : %s and %s", paramKey, PostMethod, QueryMethod)
	}

	if err := validate.Var(providerParamValue, "required,max=25"); paramKey != "message" && err != nil {
		return fmt.Errorf("the field provider.parameters.%s.param_value is required and must be a string at most 25 characters long", paramKey)
	}

	return nil
}

// validateProviderConfig validate provider configuration
func validateProviderConfig(v *viper.Viper, validate *validator.Validate) error {
	// Validate provider url
	providerUrl := v.GetString("easy_api_prom_sms_alert.provider.url")
	if err := validate.Var(providerUrl, "required,url"); err != nil {
		return fmt.Errorf("the field provider.url is required and must be a valid url")
	}

	// validate provider authentication config
	providerAuthEnabled := v.GetBool("easy_api_prom_sms_alert.provider.authentication.enabled")
	providerAuthAuthorizationType := v.GetString("easy_api_prom_sms_alert.provider.authentication.authorization.type")
	providerAuthAuthorizationCredential := v.GetString("easy_api_prom_sms_alert.provider.authentication.authorization.credential")
	if providerAuthEnabled {
		if len(providerAuthAuthorizationType) == 0 {
			return fmt.Errorf("when provider.authentication is enabled, you should provider authorization config")
		}

		if len(providerAuthAuthorizationType) != 0 {
			if err := validate.Var(providerAuthAuthorizationType, "required,oneof=Bearer Basic ApiKey"); err != nil {
				return fmt.Errorf("when provider.authentication.authorization is used, the field provider.authentication.authorization.type must be among the values : Bearer, Basic, ApiKey")
			}
			if err := validate.Var(providerAuthAuthorizationCredential, "required,max=255"); err != nil {
				return fmt.Errorf("when provider.authentication.authorization is used, the field provider.authentication.authorization.credential is required and must be a string at most 255 characters long")
			}
		}
	} else {
		v.Set("easy_api_prom_sms_alert.provider.authentication.authorization.type", "")
		v.Set("easy_api_prom_sms_alert.provider.authentication.authorization.credential", "")
	}

	// validate provider fields config
	providerParamKeys := [3]string{"from", "to", "message"}
	for _, paramKey := range providerParamKeys {
		if err := validateProviderParameter(v, validate, paramKey); err != nil {
			return err
		}
	}

	return nil
}
