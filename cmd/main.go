package main

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/app"
	"easy-api-prom-alert-sms/logging"

	"github.com/go-playground/validator/v10"
)

func main() {
	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

	configFlag, err := config.LoadConfigFlag(validate)
	if err != nil {
		logging.Log(logging.Error, "%s", err.Error())
		return
	}

	viperInstance, err := config.ReadConfigFile(configFlag.ConfigFile)
	if err != nil {
		logging.Log(logging.Error, "%s", err.Error())
		return
	}

	configLoaded, err := config.LoadConfig(viperInstance, validate)
	if err != nil {
		logging.Log(logging.Error, "%s", err.Error())
		return
	}
	logging.Log(logging.Info, "configuration file '%s' was loaded successfully", configFlag.ConfigFile)

	app.Run(configLoaded, configFlag)
}
