package config

import (
	"easy-api-prom-alert-sms/logging"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	EasyAPIPromAlertSMS struct {
		*Auth      `mapstructure:"auth"`
		Simulation bool `mapstructure:"simulation"`
		*Provider  `mapstructure:"provider"`
		Recipients `mapstructure:"recipients"`
	} `mapstructure:"easy_api_prom_sms_alert"`
}

// SetConfigDefaults sets defaults configurations values
func setConfigDefaults(v *viper.Viper) {
	v.SetDefault("easy_api_prom_sms_alert.simulation", true)
	v.SetDefault("easy_api_prom_sms_alert.auth.enabled", false)
	v.SetDefault("easy_api_prom_sms_alert.auth.username", "")
	v.SetDefault("easy_api_prom_sms_alert.auth.password", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.url", "http://localhost:5797")
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.enabled", false)
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.basic.username", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.basic.password", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.authorization.type", "Bearer")
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.authorization.credential", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.from.param_name", "from")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.from.param_value", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.from.param_method", PostMethod)
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.to.param_name", "to")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.to.param_value", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.to.param_method", PostMethod)
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.message.param_name", "")
	v.Set("easy_api_prom_sms_alert.provider.parameters.message.param_value", "")
	v.Set("easy_api_prom_sms_alert.provider.parameters.message.param_method", PostMethod)
	v.SetDefault("easy_api_prom_sms_alert.provider.timeout", "10s")
	v.SetDefault("easy_api_prom_sms_alert.recipients", make(Recipients, 0))
}

// validateConfig validate the entire configuration
func validateConfig(v *viper.Viper, validate *validator.Validate) error {
	if err := validateAuthConfig(v, validate); err != nil {
		return err
	}

	if err := validateProviderConfig(v, validate); err != nil {
		return err
	}

	if err := validateRecipientsConfig(v, validate); err != nil {
		return err
	}

	return nil
}

// LoadConfig load yaml configuration file
func LoadConfig(filename string, validate *validator.Validate) (*Config, error) {
	// Load configuration file
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logging.Log(logging.Error, err.Error())
			return nil, err
		} else {
			logging.Log(logging.Error, err.Error())
			return nil, err
		}
	}

	// Set defaut configuration
	setConfigDefaults(viper.GetViper())

	// Validate configuration file
	if err := validateConfig(viper.GetViper(), validate); err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}

	// Parse configuration file to Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}

	return &config, nil
}
