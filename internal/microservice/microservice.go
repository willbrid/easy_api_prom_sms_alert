package microservice

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/microservice/alert"
)

type IAlertMicroservice interface {
	Consume(url string, body string) error
	GetUrlAndBody(dest string, message string) (string, string, error)
}

type Microservice struct {
	IAlertMicroservice IAlertMicroservice
}

func NewMicroservice(provider *config.Provider) *Microservice {
	return &Microservice{
		IAlertMicroservice: alert.NewAlertMicroservice(provider),
	}
}
