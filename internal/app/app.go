package app

import (
	"easy-api-prom-alert-sms/alert"
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/logging"

	"net/http"

	"github.com/gorilla/mux"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag) {
	var err error

	alertSender := alert.NewAlertSender(cfgfile)
	router := mux.NewRouter()
	router.HandleFunc("/api-alert", alertSender.AlertHandler).Methods("POST")
	router.Use(alertSender.AuthMiddleware)

	logging.Log(logging.Info, "server is listening on port %v", cfgflag.ListenPort)
	if cfgflag.EnableHttps {
		logging.Log(logging.Info, "server is using https")
		err = http.ListenAndServeTLS(":"+cfgflag.ListenPort, cfgflag.CertFile, cfgflag.KeyFile, router)
	} else {
		logging.Log(logging.Info, "server is using http")
		err = http.ListenAndServe(":"+cfgflag.ListenPort, router)
	}

	if err != nil {
		logging.Log(logging.Error, "failed to start server: %v", err.Error())
	}
}
