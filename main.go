package main

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"

	"flag"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func main() {
	// Load configuration file
	var configFile string
	flag.StringVar(&configFile, "config-file", "config.yaml", "Chemin du fichier de configuration")
	flag.Parse()

	validate = validator.New(validator.WithRequiredStructEnabled())
	_, err := config.LoadConfig(configFile, validate)
	if err != nil {
		logging.Log(logging.Error, "error loading configuration")
		return
	}

	// The logic of sms alert
}
