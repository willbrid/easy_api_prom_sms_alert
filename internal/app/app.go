package app

import (
	"easy-api-prom-alert-sms/alert"
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/pkg/logger"

	"net/http"

	"github.com/gorilla/mux"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag, l *logger.Logger) {
	var err error

	alertSender := alert.NewAlertSender(cfgfile)
	router := mux.NewRouter()
	router.HandleFunc("/api-alert", alertSender.AlertHandler).Methods("POST")
	router.Use(alertSender.AuthMiddleware)

	l.Info("server is listening on port %v", cfgflag.ListenPort)
	if cfgflag.EnableHttps {
		l.Info("server is using https")
		err = http.ListenAndServeTLS(":"+cfgflag.ListenPort, cfgflag.CertFile, cfgflag.KeyFile, router)
	} else {
		l.Info("server is using http")
		err = http.ListenAndServe(":"+cfgflag.ListenPort, router)
	}

	if err != nil {
		l.Error("failed to start server: %v", err.Error())
	}
}
