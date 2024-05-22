package main

import (
	"easy-api-prom-alert-sms/alert"
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"

	"flag"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

func main() {
	var (
		configFile string
		listenPort int
	)

	flag.StringVar(&configFile, "config-file", "config.default.yaml", "Chemin du fichier de configuration")
	flag.IntVar(&listenPort, "port", 5957, "port d'écoute")
	flag.Parse()
	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Var(listenPort, "required,min=1024,max=49151"); err != nil {
		logging.Log(logging.Error, "you should provide a port number between 1024 and 49151")
		return
	}

	configLoaded, err := config.LoadConfig(configFile, validate)
	if err != nil {
		logging.Log(logging.Error, "error loading configuration")
		return
	}
	logging.Log(logging.Info, "configuration file '%s' was loaded successfully", configFile)

	alertSender := alert.NewAlertSender(configLoaded)
	router := mux.NewRouter()
	router.HandleFunc("/api-alert", alertSender.AlertHandler).Methods("POST")

	logging.Log(logging.Info, "server is listening on port %v", strconv.Itoa(listenPort))
	err = http.ListenAndServe(":"+strconv.Itoa(listenPort), router)
	if err != nil {
		logging.Log(logging.Error, "failed to start server: %v", err.Error())
		return
	}
}
