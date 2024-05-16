package config

import (
	"easy-api-prom-alert-sms/logging"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// validateAuthConfig validate authentication configuration
func validateAuthConfig(v *viper.Viper, validate *validator.Validate) error {
	authEnabled := v.GetBool("easy_api_prom_sms_alert.auth.enabled")
	authUsername := v.GetString("easy_api_prom_sms_alert.auth.username")
	authPassword := v.GetString("easy_api_prom_sms_alert.auth.password")

	if authEnabled {
		if err := validate.Var(authUsername, "required,min=2,max=25"); err != nil {
			logging.Log(logging.Error, "failed to validate auth.username field")
			return fmt.Errorf("the field auth.username is required and must be a string between 2 and 25 characters long")
		}

		if err := validate.Var(authPassword, "required,min=8,max=255"); err != nil {
			logging.Log(logging.Error, "failed to validate auth.password")
			return fmt.Errorf("the field auth.password is required and must be a string between 8 and 255 characters long")
		}

		logging.Log(logging.Info, "auth config validation passed successfully")
		return nil
	}

	v.Set("easy_api_prom_sms_alert.auth.username", "")
	v.Set("easy_api_prom_sms_alert.auth.password", "")

	logging.Log(logging.Info, "auth config validation passed successfully")
	return nil
}

func validateConfig(v *viper.Viper, validate *validator.Validate) error {
	// Validate authentication config
	if err := validateAuthConfig(v, validate); err != nil {
		return err
	}

	return nil
}
