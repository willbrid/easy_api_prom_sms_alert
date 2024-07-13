package config

import (
	"easy-api-prom-alert-sms/logging"
	"easy-api-prom-alert-sms/utils/httpclient"

	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Auth struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username" validate:"required_if=Enabled true,min=2,max=25"`
	Password string `mapstructure:"password" validate:"required_if=Enabled true,min=8"`
}

type Parameter struct {
	ParamName   string `mapstructure:"param_name" validate:"required,max=25"`
	ParamValue  string `mapstructure:"param_value" validate:"required"`
	ParamMethod string `mapstructure:"param_method" validate:"required,oneof=post query"`
}

type Provider struct {
	Url            string        `mapstructure:"url" validate:"required,url"`
	Timeout        time.Duration `mapstructure:"timeout" validate:"required"`
	ContentType    string        `mapstructure:"content_type" validate:"required,oneof=application/json application/x-www-form-urlencoded"`
	Authentication struct {
		Enabled                 bool   `mapstructure:"enabled"`
		AuthorizationType       string `mapstructure:"authorization_type" validate:"required_if=Enabled true,max=25"`
		AuthorizationCredential string `mapstructure:"authorization_credential" validate:"required_if=Enabled true"`
	} `mapstructure:"authentication"`
	Parameters struct {
		From        Parameter   `mapstructure:"from" validate:"required"`
		To          Parameter   `mapstructure:"to" validate:"required"`
		Message     Parameter   `mapstructure:"message" validate:"required"`
		ExtraParams []Parameter `mapstructure:"extra_params"`
	} `mapstructure:"parameters" validate:"required"`
}

type Recipient struct {
	Name    string   `mapstructure:"name" validate:"required,max=25"`
	Members []string `mapstructure:"members" validate:"gt=0,required,dive,min=1,max=25"`
}

type Config struct {
	EasyAPIPromAlertSMS struct {
		Auth       `mapstructure:"auth"`
		Simulation bool `mapstructure:"simulation"`
		Provider   `mapstructure:"provider"`
		Recipients []Recipient `mapstructure:"recipients" validate:"gt=0,required,dive"`
	} `mapstructure:"easy_api_prom_sms_alert"`
}

const AlertMessageTemplate string = "<%message%>"

// SetConfigDefaults sets defaults configurations values
func setConfigDefaults(v *viper.Viper) {
	v.SetDefault("easy_api_prom_sms_alert.simulation", true)
	v.SetDefault("easy_api_prom_sms_alert.auth.enabled", false)
	v.SetDefault("easy_api_prom_sms_alert.auth.username", "")
	v.SetDefault("easy_api_prom_sms_alert.auth.password", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.url", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.content_type", "application/json")
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.enabled", false)
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.authorization_type", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.authentication.authorization_credential", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.from.param_name", "from")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.from.param_value", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.from.param_method", httpclient.PostMethod)
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.to.param_name", "to")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.to.param_value", "")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.to.param_method", httpclient.PostMethod)
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.message.param_name", "content")
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.message.param_value", AlertMessageTemplate)
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.message.param_method", httpclient.PostMethod)
	v.SetDefault("easy_api_prom_sms_alert.provider.parameters.extra_params", make([]Parameter, 0))
	v.SetDefault("easy_api_prom_sms_alert.provider.timeout", "10s")
	v.SetDefault("easy_api_prom_sms_alert.recipients", make([]Recipient, 0))
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

	var viperInstance *viper.Viper = viper.GetViper()
	// Set defaut configuration
	setConfigDefaults(viperInstance)

	// Parse configuration file to Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logging.Log(logging.Error, err.Error())
		return nil, err
	}

	// Validate config struct
	if err := validate.Struct(config); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf("validation failed on field '%s' for condition '%s'", err.Field(), err.Tag())
		}
	}

	return &config, nil
}
