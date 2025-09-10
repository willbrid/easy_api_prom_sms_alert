package app

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/handler"
	ams "easy-api-prom-alert-sms/internal/microservice/alert"
	auc "easy-api-prom-alert-sms/internal/usecase/alert"
	"easy-api-prom-alert-sms/pkg/httpserver"
	"easy-api-prom-alert-sms/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	logger *logger.Logger
}

func NewApp(l *logger.Logger) *App {
	return &App{logger: l}
}

func (app *App) Run(cfgfile *config.Config, cfgflag *config.ConfigFlag) {
	alertUserCase := auc.NewAlertUseCase(
		ams.NewAlertMicroservice(cfgfile.EasyAPIPromAlertSMS.Provider),
		cfgfile.EasyAPIPromAlertSMS.Recipients,
		cfgfile.EasyAPIPromAlertSMS.Provider.Parameters.To.ParamValue,
		cfgfile.EasyAPIPromAlertSMS.Simulation,
	)

	httpServer := httpserver.NewServer(
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	handler.NewRouter(httpServer.Router, cfgfile, alertUserCase, app.logger)
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
