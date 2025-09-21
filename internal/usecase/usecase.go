package usecase

import (
	"easy-api-prom-alert-sms/internal/domain"
	"easy-api-prom-alert-sms/internal/microservice"
	"easy-api-prom-alert-sms/internal/usecase/alert"
)

type IAlertUsecase interface {
	Send(domain.Alert) error
}

type Usecases struct {
	IAlertUsecase IAlertUsecase
}

type Deps struct {
	Microservice *microservice.Microservice
	AlertConfig  *domain.AlertConfig
}

func NewUsecases(deps *Deps) *Usecases {
	return &Usecases{
		IAlertUsecase: alert.NewAlertUseCase(deps.Microservice.IAlertMicroservice, deps.AlertConfig),
	}
}
