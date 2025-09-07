package app

import (
	"easy-api-prom-alert-sms/alert"
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/pkg/httpserver"
	"easy-api-prom-alert-sms/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

type App struct {
	logger *logger.Logger
}

func NewApp(l *logger.Logger) *App {
	return &App{logger: l}
}

func (app *App) Run(cfgfile *config.Config, cfgflag *config.ConfigFlag) {
	alertSender := alert.NewAlertSender(cfgfile)
	router := mux.NewRouter()
	router.HandleFunc("/api-alert", alertSender.AlertHandler).Methods("POST")
	router.Use(alertSender.AuthMiddleware)

	httpServer := httpserver.NewServer(
		router,
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	httpServer.Start()
	var logInfoServer string
	if cfgflag.EnableHttps {
		logInfoServer = "app server is listening on port %v using https"
	} else {
		logInfoServer = "app server is listening on port %v using http"
	}
	app.logger.Info(logInfoServer, cfgflag.ListenPort)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		app.logger.Info("app server - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		app.logger.Error("app server error: %v", err.Error())
	}

	if err := httpServer.Stop(); err != nil {
		app.logger.Error("app server - stop - error: %v", err)
	}
}
