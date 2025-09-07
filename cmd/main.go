package main

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/app"
	"easy-api-prom-alert-sms/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func main() {
	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	l := logger.NewLogger()

	configFlag, err := config.LoadConfigFlag(validate)
	if err != nil {
		l.Fatal("failed to load configuration flags: %v", err.Error())
	}

	viperInstance, err := config.ReadConfigFile(configFlag.ConfigFile)
	if err != nil {
		l.Fatal("failed to read configuration file: %v", err.Error())
	}

	configLoaded, err := config.LoadConfig(viperInstance, validate)
	if err != nil {
		l.Fatal("failed to load configuration file: %v", err.Error())
	}

	l.Info("configuration file '%s' was loaded successfully", configFlag.ConfigFile)

	appInstance := app.NewApp(l)
	appInstance.Run(configLoaded, configFlag)
}
