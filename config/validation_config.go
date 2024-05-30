package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// validateConfig validate the entire configuration
func validateConfig(v *viper.Viper, validate *validator.Validate) error {
	// Validate authentication config
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
