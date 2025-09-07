package app

import (
	"easy-api-prom-alert-sms/alert"
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/pkg/logger"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	logger *logger.Logger
}

func NewApp(l *logger.Logger) *App {
	return &App{logger: l}
}

func (app *App) Run(cfgfile *config.Config, cfgflag *config.ConfigFlag) {
	var err error

	alertSender := alert.NewAlertSender(cfgfile)
	router := mux.NewRouter()
	router.HandleFunc("/api-alert", alertSender.AlertHandler).Methods("POST")
	router.Use(alertSender.AuthMiddleware)

	app.logger.Info("server is listening on port %v", cfgflag.ListenPort)
	if cfgflag.EnableHttps {
		app.logger.Info("server is using https")
		err = http.ListenAndServeTLS(":"+fmt.Sprint(cfgflag.ListenPort), cfgflag.CertFile, cfgflag.KeyFile, router)
	} else {
		app.logger.Info("server is using http")
		err = http.ListenAndServe(":"+fmt.Sprint(cfgflag.ListenPort), router)
	}

	if err != nil {
		app.logger.Error("failed to start server: %v", err.Error())
	}
}
