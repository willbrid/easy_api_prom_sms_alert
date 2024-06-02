package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Recipients []Recipient

type Recipient struct {
	Name    string   `mapstructure:"name"`
	Members []string `mapstructure:"members"`
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

		name := recipientMap["name"].(string)
		if err := validate.Var(name, "required,max=25"); err != nil {
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
			memberStr := member.(string)
			if err := validate.Var(memberStr, "required,max=25"); err != nil {
				return fmt.Errorf("the field recipients[].members[] is required and must be a string at most 25 characters long")
			}
		}
	}

	return nil
}
