package microservice

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/microservice/alert"
)

type IAlertMicroservice interface {
	Consume(url string, body string) error
	GetUrlAndBody(dest string, message string) (string, string, error)
}

type Microservices struct {
	IAlertMicroservice IAlertMicroservice
}

func NewMicroservices(provider *config.Provider) *Microservices {
	return &Microservices{
		IAlertMicroservice: alert.NewAlertMicroservice(provider),
	}
}
