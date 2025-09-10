package usecase

import (
	"easy-api-prom-alert-sms/internal/entity"
	"easy-api-prom-alert-sms/pkg/logger"
)

type IAlert interface {
	Send(entity.Alert, logger.ILogger) error
}
