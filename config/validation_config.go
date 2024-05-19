package config

import (
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
			return fmt.Errorf("the field auth.username is required and must be a string between 2 and 25 characters long")
		}

		if err := validate.Var(authPassword, "required,min=8,max=255"); err != nil {
			return fmt.Errorf("the field auth.password is required and must be a string between 8 and 255 characters long")
		}

		return nil
	}

	v.Set("easy_api_prom_sms_alert.auth.username", "")
	v.Set("easy_api_prom_sms_alert.auth.password", "")

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
	providerAuthBasicUsername := v.GetString("easy_api_prom_sms_alert.provider.authentication.basic.username")
	providerAuthBasicPassword := v.GetString("easy_api_prom_sms_alert.provider.authentication.basic.password")
	providerAuthAuthorizationType := v.GetString("easy_api_prom_sms_alert.provider.authentication.authorization.type")
	providerAuthAuthorizationCredential := v.GetString("easy_api_prom_sms_alert.provider.authentication.authorization.credential")
	if providerAuthEnabled {
		if (len(providerAuthBasicUsername) == 0 && len(providerAuthAuthorizationType) == 0) || (len(providerAuthBasicUsername) != 0 && len(providerAuthAuthorizationType) != 0) {
			return fmt.Errorf("when provider.authentication is enabled, you should provider either basic or authorization config but not both")
		}

		if len(providerAuthBasicUsername) != 0 {
			v.Set("easy_api_prom_sms_alert.provider.authentication.authorization.type", "")
			v.Set("easy_api_prom_sms_alert.provider.authentication.authorization.credential", "")

			if err := validate.Var(providerAuthBasicUsername, "required,max=255"); err != nil {
				return fmt.Errorf("when provider.authentication.basic is used, the field provider.authentication.basic.username is required and must be a string at most 255 characters long")
			}
			if err := validate.Var(providerAuthBasicPassword, "required,max=255"); err != nil {
				return fmt.Errorf("when provider.authentication.basic is used, the field provider.authentication.basic.password is required and must be a string at most 255 characters long")
			}
		}

		if len(providerAuthAuthorizationType) != 0 {
			v.Set("easy_api_prom_sms_alert.provider.authentication.basic.username", "")
			v.Set("easy_api_prom_sms_alert.provider.authentication.basic.password", "")

			if err := validate.Var(providerAuthAuthorizationType, "required,oneof='Bearer Basic ApiKey'"); err != nil {
				return fmt.Errorf("when provider.authentication.authorization is used, the field provider.authentication.authorization.type must be among the values : Bearer, Basic, ApiKey")
			}
			if err := validate.Var(providerAuthAuthorizationCredential, "required,max=255"); err != nil {
				return fmt.Errorf("when provider.authentication.authorization is used, the field provider.authentication.authorization.credential is required and must be a string at most 255 characters long")
			}
		}
	} else {
		v.Set("easy_api_prom_sms_alert.provider.authentication.basic.username", "")
		v.Set("easy_api_prom_sms_alert.provider.authentication.basic.password", "")
		v.Set("easy_api_prom_sms_alert.provider.authentication.authorization.type", "")
		v.Set("easy_api_prom_sms_alert.provider.authentication.authorization.credential", "")
	}

	// validate provider fields config
	providerFieldsFrom := v.GetString("easy_api_prom_sms_alert.provider.fields.from")
	providerFieldsFromValue := v.GetString("easy_api_prom_sms_alert.provider.fields.from_value")
	providerFieldsTo := v.GetString("easy_api_prom_sms_alert.provider.fields.to")
	providerFieldsMessage := v.GetString("easy_api_prom_sms_alert.provider.fields.message")
	if err := validate.Var(providerFieldsFrom, "required,max=25"); err != nil {
		return fmt.Errorf("the field provider.fields.from is required and must be a string at most 25 characters long")
	}
	if err := validate.Var(providerFieldsFromValue, "required,max=25"); err != nil {
		return fmt.Errorf("the field provider.fields.from_value is required and must be a string at most 25 characters long")
	}
	if err := validate.Var(providerFieldsTo, "required,max=25"); err != nil {
		return fmt.Errorf("the field provider.fields.to is required and must be a string at most 25 characters long")
	}
	if err := validate.Var(providerFieldsMessage, "required,max=25"); err != nil {
		return fmt.Errorf("the field provider.fields.message is required and must be a string at most 25 characters long")
	}

	return nil
}

func validateRecipientsConfig(v *viper.Viper, validate *validator.Validate) error {
	recipients := v.Get("easy_api_prom_sms_alert.recipients")
	recipientList, ok := recipients.([]interface{})
	if !ok {
		return fmt.Errorf("error converting recipients to slice of interface{}")
	}

	if len(recipientList) > 50 {
		return fmt.Errorf("recipients configuration must contain at most 50 items")
	}

	for _, recipient := range recipientList {
		recipientMap, ok := recipient.(map[string]interface{})
		if !ok {
			return fmt.Errorf("error converting recipient to map[string]interface{}")
		}

		if err := validate.Var(recipientMap["name"], "required,max=25"); err != nil {
			return fmt.Errorf("the field recipients[].name is required and must be a string at most 25 characters long")
		}

		members, ok := recipientMap["members"].([]interface{})
		if !ok {
			return fmt.Errorf("error converting recipientMap['members'] to slice of interface{}")
		}
		if len(members) > 50 {
			return fmt.Errorf("recipients members configuration must contain at most 50 items")
		}

		for _, member := range members {
			memberMap, ok := member.(map[string]interface{})
			if !ok {
				return fmt.Errorf("error converting member to map[string]interface{}")
			}

			if err := validate.Var(memberMap["name"], "required,max=25"); err != nil {
				return fmt.Errorf("the field recipients[].members[].name is required and must be a string at most 25 characters long")
			}
			if err := validate.Var(memberMap["phone"], "required,max=25"); err != nil {
				return fmt.Errorf("the field recipients[].members[].phone is required and must be a string at most 25 characters long")
			}
		}
	}

	return nil
}

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
