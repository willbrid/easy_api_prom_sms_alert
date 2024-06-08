package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Auth struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username" validate:"required_with=Enabled,min=2,max=25"`
	Password string `mapstructure:"password" validate:"required_with=Enabled,min=8"`
}

const (
	errParsedAuth   string = "unable to parse easy_api_prom_sms_alert.auth"
	errAuthUsername string = "the field auth.username is required and must be a string between 2 and 25 characters long"
	errAuthPassword string = "the field auth.password is required and must must have at least 8 characters long"
)

// validateAuthConfig validate authentication configuration
func validateAuthConfig(v *viper.Viper, validate *validator.Validate) error {
	authField := "easy_api_prom_sms_alert.auth"
	var auth Auth

	if err := v.Sub(authField).Unmarshal(&auth); err != nil {
		return fmt.Errorf(errParsedAuth+": %w", err)
	}

	if auth.Enabled {
		err := validate.Struct(auth)

		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return err
			}

			for _, err := range err.(validator.ValidationErrors) {
				switch err.Field() {
				case "Username":
					return fmt.Errorf(errAuthUsername)
				case "Password":
					return fmt.Errorf(errAuthPassword)
				default:
					return fmt.Errorf("validation failed for " + authField)
				}
			}
		}

		return nil
	}

	v.Set("easy_api_prom_sms_alert.auth.username", "")
	v.Set("easy_api_prom_sms_alert.auth.password", "")

	return nil
}
