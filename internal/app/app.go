package app

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/domain"
	"easy-api-prom-alert-sms/internal/handler"
	"easy-api-prom-alert-sms/internal/microservice"
	"easy-api-prom-alert-sms/internal/usecase"
	"easy-api-prom-alert-sms/pkg/httpserver"
	"easy-api-prom-alert-sms/pkg/logger"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgfile *config.Config, cfgflag *config.ConfigFlag) {
	microservice := microservice.NewMicroservices(&cfgfile.EasyAPIPromAlertSMS.Provider)
	usecases := usecase.NewUsecases(&usecase.Deps{
		Microservices: microservice,
		AlertConfig: &domain.AlertConfig{
			Recipients:           cfgfile.EasyAPIPromAlertSMS.Recipients,
			DefaultRecipientName: cfgfile.EasyAPIPromAlertSMS.Provider.Parameters.To.ParamValue,
			Simulation:           cfgfile.EasyAPIPromAlertSMS.Simulation,
		},
	})

	httpServer := httpserver.NewServer(
		fmt.Sprint(":"+fmt.Sprint(cfgflag.ListenPort)),
		cfgflag.EnableHttps,
		cfgflag.CertFile,
		cfgflag.KeyFile,
	)
	handlerInstance := handler.NewHandler(usecases, httpServer.Router)
	handlerInstance.InitRouter(cfgfile)
	httpServer.Start()

	scheme := map[bool]string{true: "https", false: "http"}[cfgflag.EnableHttps]
	logger.Info("app server is listening on port %v using %s", cfgflag.ListenPort, scheme)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app server - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		logger.Error("app server error: %v", err.Error())
	}

	if err := httpServer.Stop(); err != nil {
		logger.Error("app server - stop - error: %v", err)
	}
}
