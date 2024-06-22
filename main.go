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
		configFile  string
		listenPort  int
		enableHttps bool
		certFile    string
		keyFile     string
	)

	flag.StringVar(&configFile, "config-file", "fixtures/config.default.yaml", "configuration file path")
	flag.IntVar(&listenPort, "port", 5957, "listening port")
	flag.BoolVar(&enableHttps, "enable-https", false, "configuration to enable https")
	flag.StringVar(&certFile, "cert-file", "fixtures/tls/server.crt", "certificat file path")
	flag.StringVar(&keyFile, "key-file", "fixtures/tls/server.key", "private key file path")
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
	router.Use(alertSender.AuthMiddleware)

	strListenPort := strconv.Itoa(listenPort)

	logging.Log(logging.Info, "server is listening on port %v", strListenPort)
	if enableHttps {
		logging.Log(logging.Info, "server is using https")
		err = http.ListenAndServeTLS(":"+strListenPort, certFile, keyFile, router)
	} else {
		logging.Log(logging.Info, "server is using http")
		err = http.ListenAndServe(":"+strListenPort, router)
	}

	if err != nil {
		logging.Log(logging.Error, "failed to start server: %v", err.Error())
		return
	}
}
