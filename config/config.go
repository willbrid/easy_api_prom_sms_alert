package config

import (
	"time"

	"easy-api-prom-alert-sms/logging"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Recipients []Recipient

type Recipient struct {
	GroupName string   `mapstructure:"name"`
	Members   []Member `mapstructure:"members"`
}

type Member struct {
	Name  string `mapstructure:"name"`
	Phone string `mapstructure:"phone"`
}

type Provider struct {
	Url             string        `mapstructure:"url"`
	Timeout         time.Duration `mapstructure:"timeout"`
	*Authentication `mapstructure:"authentication"`
	*Field          `mapstructure:"fields"`
}

type Authentication struct {
	Enabled        bool `mapstructure:"enabled"`
	*Basic         `mapstructure:"basic"`
	*Authorization `mapstructure:"authorization"`
}

type Basic struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Authorization struct {
	Type       string `mapstructure:"type"`
	Credential string `mapstructure:"credential"`
}

type Field struct {
	From      string `mapstructure:"from"`
	FromValue string `mapstructure:"from_value"`
	To        string `mapstructure:"to"`
	Message   string `mapstructure:"message"`
}

type Auth struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

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
	v.SetDefault("easy_api_prom_sms_alert.provider.fields.from", "from")
	v.SetDefault("easy_api_prom_sms_alert.provider.fields.from_value", "SENDER")
	v.SetDefault("easy_api_prom_sms_alert.provider.fields.to", "to")
	v.SetDefault("easy_api_prom_sms_alert.provider.fields.message", "content")
	v.SetDefault("easy_api_prom_sms_alert.provider.timeout", "10s")
	v.SetDefault("easy_api_prom_sms_alert.recipients", make(Recipients, 0))
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
