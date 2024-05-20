package main

import (
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

	_, err := config.LoadConfig(configFile, validate)
	if err != nil {
		logging.Log(logging.Error, "error loading configuration")
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/api-alert", alertHandler).Methods("POST")
	logging.Log(logging.Info, "Server is listening on port %v", strconv.Itoa(listenPort))
	err = http.ListenAndServe(":"+strconv.Itoa(listenPort), router)
	if err != nil {
		logging.Log(logging.Error, "Failed to start server: %v", err.Error())
		return
	}
}
