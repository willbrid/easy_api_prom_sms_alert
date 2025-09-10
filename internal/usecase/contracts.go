package usecase

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/entity"
)

type IAlert interface {
	Send(entity.Alert, *config.Config) error
}
